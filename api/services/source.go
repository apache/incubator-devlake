package services

import (
	"encoding/json"
	"net/http"

	"github.com/merico-dev/lake/api/errors"
	"github.com/merico-dev/lake/api/models"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

// NewSource create source for plugin
func NewSource(data types.CreateSource) (*models.Source, error) {
	b, err := json.Marshal(data.Options)
	if err != nil {
		logger.Error(err)
		return nil, errors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	source := models.Source{
		Plugin:  data.Plugin,
		Name:    data.Name,
		Options: b,
	}
	err = db.Save(&source).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return &source, err
}
