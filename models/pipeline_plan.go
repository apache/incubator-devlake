package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type PipelinePlan struct {
	common.Model
	Name     string
	CronTime string         //cron job format, like '5 * * * 2'
	Tasks    datatypes.JSON `json:"tasks"`
	Enable   string
}
