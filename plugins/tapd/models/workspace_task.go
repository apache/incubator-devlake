package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TapdWorkSpaceTask struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	TaskId       uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkSpaceTask) TableName() string {
	return "_tool_tapd_workspace_tasks"
}
