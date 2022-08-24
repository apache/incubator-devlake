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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/utils"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	v11 "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
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
		db.Model(&models.Pipeline{}).Where("status <> ?", models.TASK_COMPLETED).Update("status", models.TASK_FAILED)
		db.Model(&models.Task{}).Where("status <> ?", models.TASK_COMPLETED).Update("status", models.TASK_FAILED)
	}

	err := ReloadBlueprints(cronManager)
	if err != nil {
		panic(err)
	}

	var pipelineMaxParallel = cfg.GetInt64("PIPELINE_MAX_PARALLEL")
	if pipelineMaxParallel < 0 {
		panic(fmt.Errorf(`PIPELINE_MAX_PARALLEL should be a positive integer`))
	}
	if pipelineMaxParallel == 0 {
		globalPipelineLog.Warn(`pipelineMaxParallel=0 means pipeline will be run No Limit`)
		pipelineMaxParallel = 10000
	}
	// run pipeline with independent goroutine
	go RunPipelineInQueue(pipelineMaxParallel)
}

// CreatePipeline and return the model
func CreatePipeline(newPipeline *models.NewPipeline) (*models.Pipeline, error) {
	// create pipeline object from posted data
	pipeline := &models.Pipeline{
		Name:          newPipeline.Name,
		FinishedTasks: 0,
		Status:        models.TASK_CREATED,
		Message:       "",
		SpentSeconds:  0,
	}
	if newPipeline.BlueprintId != 0 {
		pipeline.BlueprintId = newPipeline.BlueprintId
	}

	// save pipeline to database
	err := db.Create(&pipeline).Error
	if err != nil {
		globalPipelineLog.Error("create pipline failed: %w", err)
		return nil, errors.InternalError
	}

	// create tasks accordingly
	for i := range newPipeline.Plan {
		for j := range newPipeline.Plan[i] {
			pipelineTask := newPipeline.Plan[i][j]
			newTask := &models.NewTask{
				PipelineTask: pipelineTask,
				PipelineId:   pipeline.ID,
				PipelineRow:  i + 1,
				PipelineCol:  j + 1,
			}
			_, err := CreateTask(newTask)
			if err != nil {
				globalPipelineLog.Error("create task for pipeline failed: %w", err)
				return nil, err
			}
			// sync task state back to pipeline
			pipeline.TotalTasks += 1
		}
	}
	if err != nil {
		globalPipelineLog.Error("save tasks for pipeline failed: %w", err)
		return nil, errors.InternalError
	}
	if pipeline.TotalTasks == 0 {
		return nil, fmt.Errorf("no task to run")
	}

	// update tasks state
	pipeline.Plan, err = json.Marshal(newPipeline.Plan)
	if err != nil {
		return nil, err
	}
	err = db.Model(pipeline).Updates(map[string]interface{}{
		"total_tasks": pipeline.TotalTasks,
		"plan":        pipeline.Plan,
	}).Error
	if err != nil {
		globalPipelineLog.Error("update pipline state failed: %w", err)
		return nil, errors.InternalError
	}

	return pipeline, nil
}

// GetPipelines by query
func GetPipelines(query *PipelineQuery) ([]*models.Pipeline, int64, error) {
	pipelines := make([]*models.Pipeline, 0)
	db := db.Model(pipelines).Order("id DESC")
	if query.BlueprintId != 0 {
		db = db.Where("blueprint_id = ?", query.BlueprintId)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Pending > 0 {
		db = db.Where("finished_at is null")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if query.Page > 0 && query.PageSize > 0 {
		offset := query.PageSize * (query.Page - 1)
		db = db.Limit(query.PageSize).Offset(offset)
	}
	err = db.Find(&pipelines).Error
	if err != nil {
		return nil, count, err
	}
	return pipelines, count, nil
}

// GetPipeline by id
func GetPipeline(pipelineId uint64) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	err := db.First(pipeline, pipelineId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("pipeline not found")
		}
		return nil, err
	}
	return pipeline, nil
}

// GetPipelineLogsArchivePath creates an archive for the logs of this pipeline and returns its file path
func GetPipelineLogsArchivePath(pipeline *models.Pipeline) (string, error) {
	logPath, err := getPipelineLogsPath(pipeline)
	if err != nil {
		return "", err
	}
	archive := fmt.Sprintf("%s/%s/logging.tar.gz", os.TempDir(), uuid.New())
	if err = utils.CreateArchive(archive, true, logPath); err != nil {
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
		pipeline := &models.Pipeline{}
		for {
			db.Where("status = ?", models.TASK_CREATED).
				Not(startedPipelineIds).
				Order("id ASC").Limit(1).Find(pipeline)
			if pipeline.ID != 0 {
				break
			}
			time.Sleep(time.Second)
		}
		startedPipelineIds = append(startedPipelineIds, pipeline.ID)
		go func() {
			defer sema.Release(1)
			globalPipelineLog.Info("run pipeline, %d", pipeline.ID)
			err = runPipeline(pipeline.ID)
			if err != nil {
				globalPipelineLog.Error("failed to run pipeline, %d: %v", pipeline.ID, err)
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
			runningPipelines := make([]models.Pipeline, 0)
			err := db.Find(&runningPipelines, "status = ?", models.TASK_RUNNING).Error
			if err != nil {
				panic(err)
			}
			progressDetails := make(map[uint64]*models.TaskProgressDetail)
			// check their status against temporal
			for _, rp := range runningPipelines {
				workflowId := getTemporalWorkflowId(rp.ID)
				desc, err := temporalClient.DescribeWorkflowExecution(
					context.Background(),
					workflowId,
					"",
				)
				if err != nil {
					globalPipelineLog.Error("failed to query workflow execution: %w", err)
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
								globalPipelineLog.Error("failed to get next from workflow history iterator: %w", err)
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
						globalPipelineLog.Error("failed to update db: %w", err)
					}
					continue
				}

				// check pending activity
				for _, activity := range desc.PendingActivities {
					taskId, err := getTaskIdFromActivityId(activity.ActivityId)
					if err != nil {
						globalPipelineLog.Error("unable to extract task id from activity id `%s`", activity.ActivityId)
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
					err = dc.FromPayload(lastPayload, progressDetail)
					if err != nil {
						globalPipelineLog.Error("failed to unmarshal heartbeat payload: %w", err)
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
func NotifyExternal(pipelineId uint64) error {
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
		globalPipelineLog.Error("failed to send notification: %w", err)
		return err
	}
	return nil
}

// CancelPipeline FIXME ...
func CancelPipeline(pipelineId uint64) error {
	if temporalClient != nil {
		return temporalClient.CancelWorkflow(context.Background(), getTemporalWorkflowId(pipelineId), "")
	}
	pendingTasks, count, err := GetTasks(&TaskQuery{PipelineId: pipelineId, Pending: 1, PageSize: -1})
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	for _, pendingTask := range pendingTasks {
		_ = CancelTask(pendingTask.ID)
	}
	return err
}

// getPipelineLogsPath gets the logs directory of this pipeline
func getPipelineLogsPath(pipeline *models.Pipeline) (string, error) {
	pipelineLog := getPipelineLogger(pipeline)
	path := pipelineLog.GetConfig().Path
	_, err := os.Stat(path)
	if err == nil {
		return filepath.Dir(path), nil
	}
	if os.IsNotExist(err) {
		return "", fmt.Errorf("logs for pipeline #%d not found. You may be missing the LOGGING_DIR setting: %w", pipeline.ID, err)
	}
	return "", fmt.Errorf("err validating logs path for pipeline #%d: %w", pipeline.ID, err)
}
