package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

const (
	TASK_CREATED   = "TASK_CREATED"
	TASK_RUNNING   = "TASK_RUNNING"
	TASK_COMPLETED = "TASK_COMPLETED"
	TASK_FAILED    = "TASK_FAILED"
)

type Task struct {
	common.Model
	Plugin       string         `json:"plugin" gorm:"index"`
	Options      datatypes.JSON `json:"options"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	Progress     float32        `json:"progress"`
	PipelineId   uint64         `json:"pipelineId" gorm:"index"`
	PipelineRow  int            `json:"pipelineRow"`
	PipelineCol  int            `json:"pipelineCol"`
	BeganAt      *time.Time     `json:"beganAt"`
	FinishedAt   *time.Time     `json:"finishedAt" gorm:"index"`
	SpentSeconds int            `json:"spentSeconds"`
}

type NewTask struct {
	// Plugin name
	Plugin string `json:"plugin" binding:"required"`
	// Options for the plugin task to be triggered
	Options     map[string]interface{} `json:"options"`
	PipelineId  uint64                 `json:"-"`
	PipelineRow int                    `json:"-"`
	PipelineCol int                    `json:"-"`
}
