package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"gorm.io/datatypes"
)

type Blueprint struct {
	Name       string         `json:"name" validate:"required"`
	Tasks      datatypes.JSON `json:"tasks" validate:"required"`
	Enable     bool           `json:"enable"`
	CronConfig string         `json:"cronConfig" validate:"required"`
	common.Model
}

func (Blueprint) TableName() string {
	return "_devlake_blueprints"
}
