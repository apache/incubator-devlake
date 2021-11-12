package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins"
)

const (
	TASK_CREATED   = "TASK_CREATED"
	TASK_COMPLETED = "TASK_COMPLETED"
	TASK_FAILED    = "TASK_FAILED"
)

// FIXME: don't use notification service here
// move it to controller
var notificationService *NotificationService
var runningTasks map[uint64]context.CancelFunc

type NewTask struct {
	// Plugin name
	Plugin string `json:"plugin" binding:"required"`
	// Options for the plugin task to be triggered
	Options map[string]interface{} `json:"options" binding:"required"`
}

type TaskQuery struct {
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Plugin   string `form:"plugin"`
	SourceId int64  `form:"source_id"`
}

func init() {
	var notificationEndpoint = config.V.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = config.V.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}
	// FIXME: don't cancel tasks here
	models.Db.Model(&models.Task{}).Where("status != ?", TASK_COMPLETED).Update("status", TASK_FAILED)
	runningTasks = make(map[uint64]context.CancelFunc)
}

func CreateTaskInDB(data NewTask) (*models.Task, error) {
	var sourceId int64
	if source, ok := data.Options["sourceId"].(float64); ok {
		sourceId = int64(source)
		if sourceId == 0 {
			return nil, fmt.Errorf("invalid sourceId: %d", sourceId)
		}
	}

	b, err := json.Marshal(data.Options)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		Plugin:   data.Plugin,
		Options:  b,
		Status:   TASK_CREATED,
		Message:  "",
		SourceId: sourceId,
	}
	err = models.Db.Save(&task).Error
	if err != nil {
		logger.Error("Database error", err)
		return nil, errors.InternalError
	}
	return &task, nil
}

func RunTask(task models.Task, data NewTask, taskComplete chan bool) (models.Task, error) {
	// trigger plugins
	data.Options["ID"] = task.ID
	ctx, cancel := context.WithCancel(context.Background())
	runningTasks[task.ID] = cancel
	go func() {
		logger.Info("run task ", task)
		progress := make(chan float32)
		go func() {
			err := plugins.RunPlugin(task.Plugin, data.Options, progress, ctx)
			if err != nil {
				logger.Error("Task error", err)
				task.Status = TASK_FAILED
				task.Message = err.Error()
			}
			err = models.Db.Save(&task).Error
			if err != nil {
				logger.Error("Database error", err)
			}
		}()

		for p := range progress {
			task.Progress = p
			logger.Info("running plugin progress", task)
			err := models.Db.Save(&task).Error
			if err != nil {
				logger.Error("Database error", err)
			}
		}
		task.Status = TASK_COMPLETED
		err := models.Db.Save(&task).Error
		if err != nil {
			logger.Error("Database error", err)
		}
		delete(runningTasks, task.ID)
		// TODO: send notification
		if notificationService != nil {
			err = notificationService.TaskSuccess(TaskSuccessNotification{
				TaskID:     task.ID,
				PluginName: task.Plugin,
				CreatedAt:  task.CreatedAt,
				UpdatedAt:  task.UpdatedAt,
			})
			if err != nil {
				logger.Error("Failed to send notification", err)
			}
		}
		taskComplete <- true
	}()
	return task, nil
}

func CancelTask(taskId uint64) error {
	if cancel, ok := runningTasks[taskId]; ok {
		logger.Info("cancel task ", taskId)
		cancel()
		delete(runningTasks, taskId)
	} else {
		return fmt.Errorf("unable to cancel task %v", taskId)
	}
	return nil
}

func GetPendingTasks() ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	whereClause := "progress < 1 AND status != 'TASK_COMPLETED'"
	db := models.Db.Model(&models.Task{}).Where(whereClause).Order("id DESC")
	err := db.Debug().Find(&tasks)
	if err != nil {
		return tasks, nil
	}
	return tasks, nil
}
func GetTasks(query *TaskQuery) ([]models.Task, int64, error) {
	db := models.Db.Model(&models.Task{}).Order("id DESC")
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Plugin != "" {
		db = db.Where("plugin = ?", query.Plugin)
	}
	if query.SourceId != 0 {
		db = db.Where("source_id = ?", query.SourceId)
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
	tasks := make([]models.Task, 0)
	err = db.Find(&tasks).Error
	if err != nil {
		return nil, count, err
	}
	return tasks, count, nil
}

func CreateTasksInDBFromJSON(data [][]NewTask) [][]models.Task {
	// create all the tasks in the db without running the tasks
	var tasks [][]models.Task

	for i := 0; i < len(data); i++ {
		var tasksToAppend []models.Task
		for j := 0; j < len(data[i]); j++ {
			task, _ := CreateTaskInDB(data[i][j])
			tasksToAppend = append(tasksToAppend, *task)
		}
		tasks = append(tasks, tasksToAppend)
	}

	return tasks
}

func RunAllTasks(data [][]NewTask, tasks [][]models.Task) (err error) {
	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	// for _, array := range data {
	for i := 0; i < len(data); i++ {

		taskComplete := make(chan bool)
		count := 0
		// for _, taskFromRequest := range array {
		for j := 0; j < len(data[i]); j++ {
			_, err := RunTask(tasks[i][j], data[i][j], taskComplete)
			if err != nil {
				return err
			}
		}
		for range taskComplete {
			count++
			if count == len(data[i]) {
				close(taskComplete)
			}
		}
	}
	return nil
}
