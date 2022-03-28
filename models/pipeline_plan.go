package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type PipelinePlan struct {
	common.Model
	Name       string         `json:"name"`
	Tasks      datatypes.JSON `json:"tasks"`
	Enable     bool           `json:"enable"`
	CronConfig string         `json:"cronConfig"`
}
type InputPipelinePlan struct {
	Name           string       `json:"name"`
	Tasks          [][]*NewTask `json:"tasks"`
	CronConfig     string       `json:"cronConfig"`
	Enable         bool         `json:"enable"`
	PipelinePlanId uint64
}

type EditPipelinePlan InputPipelinePlan
