package services

import (
	"encoding/json"
	"net/http"

	"github.com/merico-dev/lake/api/errors"
	"github.com/merico-dev/lake/api/models"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

func NewTask(data types.CreateTask) (*models.Task, error) {
	b, err := json.Marshal(data.Options)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		Plugin:  data.Plugin,
		Options: b,
	}
	err = db.Save(&task).Error
	if err != nil {
		logger.Error("Database error", err)
		return nil, errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return &task, nil
}
