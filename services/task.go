package services

import (
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

type NewTask struct {
	// Plugin name
	Plugin string `json:"plugin" binding:"required"`
	// Options for the plugin task to be triggered
	Options map[string]interface{} `json:"options" binding:"required"`
}

func init() {
	var notificationEndpoint = config.V.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = config.V.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}
	// FIXME: don't cancel tasks here
	models.Db.Model(&models.Task{}).Where("status != ?", TASK_COMPLETED).Update("status", TASK_FAILED)
}

func CreateTaskInDB(data NewTask) (*models.Task, error) {
	b, err := json.Marshal(data.Options)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		Plugin:  data.Plugin,
		Options: b,
		Status:  TASK_CREATED,
		Message: "",
	}
	logger.Info("JON >>> PRE: saved task in DB", task)
	err = models.Db.Save(&task).Error
	if err != nil {
		logger.Error("Database error", err)
		return nil, errors.InternalError
	}
	logger.Info("JON >>> saved task in DB", true)
	return &task, nil
}

func RunTask(data NewTask, taskComplete chan bool) (*models.Task, error) {
	task, _ := CreateTaskInDB(data)
	// trigger plugins
	data.Options["ID"] = task.ID
	go func() {
		progress := make(chan float32)
		go func() {
			err := plugins.RunPlugin(task.Plugin, data.Options, progress)
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
			fmt.Printf("running plugin %v, progress: %v\n", task.Plugin, p*100)
			task.Progress = p
			models.Db.Save(&task)
		}
		task.Status = TASK_COMPLETED
		err := models.Db.Save(&task).Error
		if err != nil {
			logger.Error("Database error", err)
		}
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

func GetTasks(status string) ([]models.Task, error) {
	db := models.Db
	if status != "" {
		db = db.Where("status = ?", status)
	}
	tasks := make([]models.Task, 0)
	err := db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
