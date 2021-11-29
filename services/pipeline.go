package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
)

var notificationService *NotificationService

type PipelineQuery struct {
	Status   string `form:"status"`
	Pending  int    `form:"pending"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

func init() {
	var notificationEndpoint = config.V.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = config.V.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}
	models.Db.Model(&models.Pipeline{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)
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

	// save pipeline to database
	err := models.Db.Create(&pipeline).Error
	if err != nil {
		logger.Error("create pipline failed", err)
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
				logger.Error("create task for pipeline failed", err)
				return nil, err
			}
			// sync task state back to pipeline
			pipeline.TotalTasks += 1
		}
	}
	if err != nil {
		logger.Error("save tasks for pipeline failed", err)
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
	err = models.Db.Model(pipeline).Updates(map[string]interface{}{
		"total_tasks": pipeline.TotalTasks,
		"tasks":       pipeline.Tasks,
	}).Error
	if err != nil {
		logger.Error("update pipline state failed", err)
		return nil, errors.InternalError
	}

	return pipeline, nil
}

func GetPipelines(query *PipelineQuery) ([]*models.Pipeline, int64, error) {
	pipelines := make([]*models.Pipeline, 0)
	db := models.Db.Model(pipelines).Order("id DESC")
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
	err := models.Db.Find(pipeline, pipelineId).Error
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func RunPipeline(pipelineId uint64) error {
	pipeline, err := GetPipeline(pipelineId)
	if err != nil {
		return err
	}
	// load tasks for pipeline
	var tasks []*models.Task
	err = models.Db.Where("pipeline_id = ?", pipeline.ID).Order("pipeline_row, pipeline_col").Find(&tasks).Error
	if err != nil {
		return err
	}
	// convert to 2d array
	taskIds := make([][]uint64, 0)
	for _, task := range tasks {
		if len(taskIds) < task.PipelineRow {
			taskIds = append(taskIds, make([]uint64, 0))
		}
		taskIds[task.PipelineRow-1] = append(taskIds[task.PipelineRow-1], task.ID)
	}

	beganAt := time.Now()
	err = models.Db.Model(pipeline).Updates(map[string]interface{}{
		"status":   models.TASK_RUNNING,
		"message":  "",
		"began_at": beganAt,
	}).Error
	if err != nil {
		return err
	}
	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	finishedTasks := 0
	rowResults := make(chan error)
	rowErrors := make([]string, 0)
	for _, row := range taskIds {
		rowFinished := 0
		for _, taskId := range row {
			taskId := taskId
			go func() {
				logger.Info("run task in background ", taskId)
				rowResults <- RunTask(taskId)
			}()
		}
		for err = range rowResults {
			finishedTasks++
			rowFinished++
			if err != nil {
				logger.Error("pipeline task failed", err)
				rowErrors = append(rowErrors, err.Error())
			}
			err = models.Db.Model(pipeline).Updates(map[string]interface{}{
				"status":         models.TASK_RUNNING,
				"finished_tasks": finishedTasks,
			}).Error
			if err != nil {
				logger.Error("update pipeline state failed", err)
				rowErrors = append(rowErrors, err.Error())
			}
			if rowFinished == len(row) {
				break
			}
		}
		if len(rowErrors) > 0 {
			err = fmt.Errorf(strings.Join(rowErrors, "\n"))
			break
		}
	}
	close(rowResults)

	logger.Info("pipeline finished:", err == nil)
	// finished, update database
	finishedAt := time.Now()
	spentSeconds := finishedAt.Unix() - beganAt.Unix()
	if err != nil {
		err = models.Db.Model(pipeline).Updates(map[string]interface{}{
			"status":        models.TASK_FAILED,
			"message":       err.Error(),
			"finished_at":   finishedAt,
			"spent_seconds": spentSeconds,
		}).Error
	} else {
		err = models.Db.Model(pipeline).Updates(map[string]interface{}{
			"status":        models.TASK_COMPLETED,
			"message":       "",
			"finished_at":   finishedAt,
			"spent_seconds": spentSeconds,
		}).Error
	}
	if err != nil {
		return err
	}

	// notify external webhook
	return NotifyExternal(pipelineId)
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
		logger.Error("Failed to send notification", err)
		return err
	}
	return nil
}

func CancelPipeline(pipelineId uint64) error {
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
