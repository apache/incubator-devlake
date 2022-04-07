package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/runner"
	"github.com/merico-dev/lake/worker/app"
	v11 "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

var notificationService *NotificationService
var temporalClient client.Client
var pipelineLog = logger.Global.Nested("pipeline service")

type PipelineQuery struct {
	Status   string `form:"status"`
	Pending  int    `form:"pending"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
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
}

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
		pipelineLog.Error("create pipline failed: %w", err)
		return nil, errors.InternalError
	}

	// create tasks accordingly
	for i := range newPipeline.Tasks {
		for j := range newPipeline.Tasks[i] {
			newTask := newPipeline.Tasks[i][j]
			newTask.PipelineId = pipeline.ID
			newTask.PipelineRow = i + 1
			newTask.PipelineCol = j + 1
			_, err := CreateTask(newTask)
			if err != nil {
				pipelineLog.Error("create task for pipeline failed: %w", err)
				return nil, err
			}
			// sync task state back to pipeline
			pipeline.TotalTasks += 1
		}
	}
	if err != nil {
		pipelineLog.Error("save tasks for pipeline failed: %w", err)
		return nil, errors.InternalError
	}
	if pipeline.TotalTasks == 0 {
		return nil, fmt.Errorf("no task to run")
	}

	// update tasks state
	pipeline.Tasks, err = json.Marshal(newPipeline.Tasks)
	if err != nil {
		return nil, err
	}
	err = db.Model(pipeline).Updates(map[string]interface{}{
		"total_tasks": pipeline.TotalTasks,
		"tasks":       pipeline.Tasks,
	}).Error
	if err != nil {
		pipelineLog.Error("update pipline state failed: %w", err)
		return nil, errors.InternalError
	}

	return pipeline, nil
}

func GetPipelines(query *PipelineQuery) ([]*models.Pipeline, int64, error) {
	pipelines := make([]*models.Pipeline, 0)
	db := db.Model(pipelines).Order("id DESC")
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

func GetPipeline(pipelineId uint64) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	err := db.Find(pipeline, pipelineId).Error
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func RunPipeline(pipelineId uint64) error {
	var err error
	// run
	if temporalClient != nil {
		err = runPipelineViaTemporal(pipelineId)
	} else {
		err = runPipelineStandalone(pipelineId)
	}
	// load
	pipeline, e := GetPipeline(pipelineId)
	if e != nil {
		return err
	}
	// finished, update database
	finishedAt := time.Now()
	pipeline.FinishedAt = &finishedAt
	pipeline.SpentSeconds = int(finishedAt.Unix() - pipeline.BeganAt.Unix())
	if err != nil {
		pipeline.Status = models.TASK_FAILED
		pipeline.Message = err.Error()
	} else {
		pipeline.Status = models.TASK_COMPLETED
		pipeline.Message = ""
	}
	dbe := db.Model(pipeline).Select("finished_at", "spent_seconds", "status", "message").Updates(pipeline).Error
	if dbe != nil {
		pipelineLog.Error("update pipeline state failed: %w", dbe)
		return dbe
	}
	// notify external webhook
	return NotifyExternal(pipelineId)
}

func getTemporalWorkflowId(pipelineId uint64) string {
	return fmt.Sprintf("pipeline #%d", pipelineId)
}

func runPipelineViaTemporal(pipelineId uint64) error {
	workflowOpts := client.StartWorkflowOptions{
		ID:        getTemporalWorkflowId(pipelineId),
		TaskQueue: cfg.GetString("TEMPORAL_TASK_QUEUE"),
	}
	// send only the very basis data
	configJson, err := json.Marshal(cfg.AllSettings())
	if err != nil {
		return err
	}
	pipelineLog.Info("enqueue pipeline #%d into temporal task queue", pipelineId)
	workflow, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOpts,
		app.DevLakePipelineWorkflow,
		configJson,
		pipelineId,
	)
	if err != nil {
		pipelineLog.Error("failed to enqueue pipeline #%d into temporal", pipelineId)
		return err
	}
	err = workflow.Get(context.Background(), nil)
	if err != nil {
		pipelineLog.Info("failed to execute pipeline #%d via temporal: %w", pipelineId, err)
	}
	pipelineLog.Info("pipeline #%d finished by temporal", pipelineId)
	return err
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
					pipelineLog.Error("failed to query workflow execution: %w", err)
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
								pipelineLog.Error("failed to get next from workflow history iterator: %w", err)
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
						pipelineLog.Error("failed to update db: %w", err)
					}
					continue
				}

				// check pending activity
				for _, activity := range desc.PendingActivities {
					taskId, err := getTaskIdFromActivityId(activity.ActivityId)
					if err != nil {
						pipelineLog.Error("unable to extract task id from activity id `%s`", activity.ActivityId)
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
						pipelineLog.Error("failed to unmarshal heartbeat payload: %w", err)
						continue
					}
				}
			}
			runningTasks.setAll(progressDetails)
		}
	}()
}

func runPipelineStandalone(pipelineId uint64) error {
	return runner.RunPipeline(
		cfg,
		pipelineLog.Nested(fmt.Sprintf("pipeline #%d", pipelineId)),
		db,
		pipelineId,
		runTasksStandalone,
	)
}

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
		pipelineLog.Error("failed to send notification: %w", err)
		return err
	}
	return nil
}

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
