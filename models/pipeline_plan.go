package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type PipelinePlan struct {
	common.Model
	Name   string         `json:"name"`
	Tasks  datatypes.JSON `json:"tasks"`
	Enable bool           `json:"enable"`
	CronConfig string
}
