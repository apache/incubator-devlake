package services

import (
	"encoding/json"

	"github.com/merico-dev/lake/api/models"
	"github.com/merico-dev/lake/api/types"
)

// NewSource create source for plugin
func NewSource(source types.CreateSource) error {
	b, err := json.Marshal(source.Options)
	if err != nil {
		return err
	}
	return db.Save(&models.Source{
		Plugin:  source.Plugin,
		Name:    source.Name,
		Options: b,
	}).Error
}
