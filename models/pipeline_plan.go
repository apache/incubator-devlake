package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type PipelinePlan struct {
	common.Model
	Name   string
	Tasks  datatypes.JSON `json:"tasks"`
	Enable bool
	CronConfig
}

type CronConfig struct {
	Type   string //weekly, monthly, interval
	Day    int
	Hour   int
	Minute int
}
