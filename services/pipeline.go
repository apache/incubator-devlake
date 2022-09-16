/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/utils"
	"github.com/google/uuid"
	v11 "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"golang.org/x/sync/semaphore"
)

var notificationService *NotificationService
var temporalClient client.Client
var globalPipelineLog = logger.Global.Nested("pipeline service")

// PipelineQuery FIXME ...
type PipelineQuery struct {
	Status      string `form:"status"`
	Pending     int    `form:"pending"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
	BlueprintId uint64 `form:"blueprint_id"`
}

func pipelineServiceInit() {
	// notification
	var notificationEndpoint = cfg.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = cfg.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}

	// temporal client
	var temporalUrl = cfg.GetString("TEMPORAL_URL")
	if temporalUrl != "" {
		// TODO: logger
		var err error
		temporalClient, err = client.NewClient(client.Options{
			HostPort: temporalUrl,
		})
		if err != nil {
			panic(err)
		}
		watchTemporalPipelines()
	} else {
		// standalone mode: reset pipeline status
		db.Model(&models.DbPipeline{}).Where("status <> ?", models.TASK_COMPLETED).Update("status", models.TASK_FAILED)
		db.Model(&models.Task{}).Where("status <> ?", models.TASK_COMPLETED).Update("status", models.TASK_FAILED)
	}

	err := ReloadBlueprints(cronManager)
	if err != nil {
		panic(err)
	}

	var pipelineMaxParallel = cfg.GetInt64("PIPELINE_MAX_PARALLEL")
	if pipelineMaxParallel < 0 {
		panic(errors.BadInput.New(`PIPELINE_MAX_PARALLEL should be a positive integer`))
	}
	if pipelineMaxParallel == 0 {
		globalPipelineLog.Warn(nil, `pipelineMaxParallel=0 means pipeline will be run No Limit`)
		pipelineMaxParallel = 10000
	}
	// run pipeline with independent goroutine
	go RunPipelineInQueue(pipelineMaxParallel)
}

// CreatePipeline and return the model
func CreatePipeline(newPipeline *models.NewPipeline) (*models.Pipeline, errors.Error) {
	dbPipeline, err := CreateDbPipeline(newPipeline)
	if err != nil {
		return nil, errors.Convert(err)
	}
	dbPipeline, err = decryptDbPipeline(dbPipeline)
	if err != nil {
		return nil, err
	}
	pipeline := parsePipeline(dbPipeline)
	return pipeline, nil
}

// GetPipelines by query
func GetPipelines(query *PipelineQuery) ([]*models.Pipeline, int64, errors.Error) {
	dbPipelines, i, err := GetDbPipelines(query)
	if err != nil {
		return nil, 0, errors.Convert(err)
	}
	pipelines := make([]*models.Pipeline, 0)
	for _, dbPipeline := range dbPipelines {
		dbPipeline, err = decryptDbPipeline(dbPipeline)
		if err != nil {
			return nil, 0, err
		}
		pipeline := parsePipeline(dbPipeline)
		pipelines = append(pipelines, pipeline)
	}

	return pipelines, i, nil
}

// GetPipeline by id
func GetPipeline(pipelineId uint64) (*models.Pipeline, errors.Error) {
	dbPipeline, err := GetDbPipeline(pipelineId)
	if err != nil {
		return nil, err
	}
	dbPipeline, err = decryptDbPipeline(dbPipeline)
	if err != nil {
		return nil, err
	}
	pipeline := parsePipeline(dbPipeline)
	return pipeline, nil
}

// GetPipelineLogsArchivePath creates an archive for the logs of this pipeline and returns its file path
func GetPipelineLogsArchivePath(pipeline *models.Pipeline) (string, errors.Error) {
	logPath, err := getPipelineLogsPath(pipeline)
	if err != nil {
		return "", err
	}
	archive := fmt.Sprintf("%s/%s/logging.tar.gz", os.TempDir(), uuid.New())
	if err = utils.CreateGZipArchive(archive, fmt.Sprintf("%s/*", logPath)); err != nil {
		return "", err
	}
	return archive, err
}

// RunPipelineInQueue query pipeline from db and run it in a queue
func RunPipelineInQueue(pipelineMaxParallel int64) {
	sema := semaphore.NewWeighted(pipelineMaxParallel)
	startedPipelineIds := []uint64{}
	for {
		globalPipelineLog.Info("wait for new pipeline")
		// start goroutine when sema lock ready and pipeline exist.
		// to avoid read old pipeline, acquire lock before read exist pipeline
		err := sema.Acquire(context.TODO(), 1)
		if err != nil {
			panic(err)
		}
		globalPipelineLog.Info("get lock and wait pipeline")
		dbPipeline := &models.DbPipeline{}
		for {
			db.Where("status = ?", models.TASK_CREATED).
				Not(startedPipelineIds).
				Order("id ASC").Limit(1).Find(dbPipeline)
			if dbPipeline.ID != 0 {
				break
			}
			time.Sleep(time.Second)
		}
		startedPipelineIds = append(startedPipelineIds, dbPipeline.ID)
		go func() {
			defer sema.Release(1)
			globalPipelineLog.Info("run pipeline, %d", dbPipeline.ID)
			err = runPipeline(dbPipeline.ID)
			if err != nil {
				globalPipelineLog.Error(err, "failed to run pipeline %d", dbPipeline.ID)
			}
		}()
	}
}

func watchTemporalPipelines() {
	ticker := time.NewTicker(3 * time.Second)
	dc := converter.GetDefaultDataConverter()
	go func() {
		// run forever
		for range ticker.C {
			// load all running pipeline from database
			runningDbPipelines := make([]models.DbPipeline, 0)
			err := db.Find(&runningDbPipelines, "status = ?", models.TASK_RUNNING).Error
			if err != nil {
				panic(err)
			}
			// progressDetails will be only used in this goroutine now
			// So it needn't lock and unlock now
			progressDetails := make(map[uint64]*models.TaskProgressDetail)
			// check their status against temporal
			for _, rp := range runningDbPipelines {
				workflowId := getTemporalWorkflowId(rp.ID)
				desc, err := temporalClient.DescribeWorkflowExecution(
					context.Background(),
					workflowId,
					"",
				)
				if err != nil {
					globalPipelineLog.Error(err, "failed to query workflow execution: %w", err)
					continue
				}
				// workflow is terminated by outsider
				s := desc.WorkflowExecutionInfo.Status
				if s != v11.WORKFLOW_EXECUTION_STATUS_RUNNING {
					rp.Status = models.TASK_COMPLETED
					if s != v11.WORKFLOW_EXECUTION_STATUS_COMPLETED {
						rp.Status = models.TASK_FAILED
						// get error message
						hisIter := temporalClient.GetWorkflowHistory(
							context.Background(),
							workflowId,
							"",
							false,
							v11.HISTORY_EVENT_FILTER_TYPE_CLOSE_EVENT,
						)
						for hisIter.HasNext() {
							his, err := hisIter.Next()
							if err != nil {
								globalPipelineLog.Error(err, "failed to get next from workflow history iterator: %w", err)
								continue
							}
							rp.Message = fmt.Sprintf("temporal event type: %v", his.GetEventType())
						}
					}
					rp.FinishedAt = desc.WorkflowExecutionInfo.CloseTime
					err = db.Model(rp).Updates(map[string]interface{}{
						"status":      rp.Status,
						"message":     rp.Message,
						"finished_at": rp.FinishedAt,
					}).Error
					if err != nil {
						globalPipelineLog.Error(err, "failed to update db: %w", err)
					}
					continue
				}

				// check pending activity
				for _, activity := range desc.PendingActivities {
					taskId, err := getTaskIdFromActivityId(activity.ActivityId)
					if err != nil {
						globalPipelineLog.Error(err, "unable to extract task id from activity id `%s`", activity.ActivityId)
						continue
					}
					progressDetail := &models.TaskProgressDetail{}
					progressDetails[taskId] = progressDetail
					heartbeats := activity.GetHeartbeatDetails()
					if heartbeats == nil {
						continue
					}
					payloads := heartbeats.GetPayloads()
					if len(payloads) == 0 {
						return
					}
					lastPayload := payloads[len(payloads)-1]
					if err := dc.FromPayload(lastPayload, progressDetail); err != nil {
						globalPipelineLog.Error(err, "failed to unmarshal heartbeat payload: %w", err)
						continue
					}
				}
			}
			runningTasks.setAll(progressDetails)
		}
	}()
}

func getTemporalWorkflowId(pipelineId uint64) string {
	return fmt.Sprintf("pipeline #%d", pipelineId)
}

// NotifyExternal FIXME ...
func NotifyExternal(pipelineId uint64) errors.Error {
	if notificationService == nil {
		return nil
	}
	// send notification to an external web endpoint
	pipeline, err := GetPipeline(pipelineId)
	if err != nil {
		return err
	}
	err = notificationService.PipelineStatusChanged(PipelineNotification{
		PipelineID: pipeline.ID,
		CreatedAt:  pipeline.CreatedAt,
		UpdatedAt:  pipeline.UpdatedAt,
		BeganAt:    pipeline.BeganAt,
		FinishedAt: pipeline.FinishedAt,
		Status:     pipeline.Status,
	})
	if err != nil {
		globalPipelineLog.Error(err, "failed to send notification: %w", err)
		return err
	}
	return nil
}

// CancelPipeline FIXME ...
func CancelPipeline(pipelineId uint64) errors.Error {
	if temporalClient != nil {
		return errors.Convert(temporalClient.CancelWorkflow(context.Background(), getTemporalWorkflowId(pipelineId), ""))
	}
	pendingTasks, count, err := GetTasks(&TaskQuery{PipelineId: pipelineId, Pending: 1, PageSize: -1})
	if err != nil {
		return errors.Convert(err)
	}
	if count == 0 {
		return nil
	}
	for _, pendingTask := range pendingTasks {
		_ = CancelTask(pendingTask.ID)
	}
	return errors.Convert(err)
}

// getPipelineLogsPath gets the logs directory of this pipeline
func getPipelineLogsPath(pipeline *models.Pipeline) (string, errors.Error) {
	pipelineLog := getPipelineLogger(pipeline)
	path := filepath.Dir(pipelineLog.GetConfig().Path)
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}
	if os.IsNotExist(err) {
		return "", errors.NotFound.Wrap(err, fmt.Sprintf("logs for pipeline #%d not found", pipeline.ID))
	}
	return "", errors.Default.Wrap(err, fmt.Sprintf("error validating logs path for pipeline #%d", pipeline.ID))
}
