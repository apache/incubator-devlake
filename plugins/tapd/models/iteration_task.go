package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdIterationTask struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	IterationId  uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`

	TaskId          uint64 `gorm:"primaryKey"`
	ResolutionDate  *helper.CSTTime
	TaskCreatedDate *helper.CSTTime
}

func (TapdIterationTask) TableName() string {
	return "_tool_tapd_iteration_tasks"
}
