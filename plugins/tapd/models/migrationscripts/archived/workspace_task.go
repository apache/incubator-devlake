package archived

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdWorkSpaceTask struct {
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceID uint64 `gorm:"primaryKey"`
	TaskId      uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkSpaceTask) TableName() string {
	return "_tool_tapd_workspace_tasks"
}
