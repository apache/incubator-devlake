package services

import (
	"encoding/json"
	"fmt"

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

type NewTask struct {
	// Plugin name
	Plugin string `json:"plugin" binding:"required"`
	// Options for the plugin task to be triggered
	Options map[string]interface{} `json:"options" binding:"required"`
}

func CreateTask(data NewTask) (*models.Task, error) {
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
	err = models.Db.Save(&task).Error
	if err != nil {
		logger.Error("Database error", err)
		return nil, errors.InternalError
	}

	// trigger plugins
	data.Options["ID"] = task.ID
	go func() {
		progress := make(chan float32)

		go func() {
			err = plugins.RunPlugin(task.Plugin, data.Options, progress)
			if err != nil {
				logger.Error("Task error", err)
				task.Status = TASK_FAILED
				task.Message = err.Error()
			}
			err := models.Db.Save(&task).Error
			if err != nil {
				logger.Error("Database error", err)
			}
		}()

		for p := range progress {
			fmt.Printf("running plugin %v, progress: %v\n", task.Plugin, p*100)
		}
		task.Status = TASK_COMPLETED
		err := models.Db.Save(&task).Error
		if err != nil {
			logger.Error("Database error", err)
		}
	}()
	return &task, nil
}
