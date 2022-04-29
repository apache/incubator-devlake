package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdIterationTask struct {
	common.NoPKModel
	ConnectionId    uint64 `gorm:"primaryKey"`
	IterationId     uint64 `gorm:"primaryKey"`
	TaskId          uint64 `gorm:"primaryKey"`
	WorkspaceID     uint64 `gorm:"primaryKey"`
	ResolutionDate  *core.Iso8601Time
	TaskCreatedDate *core.Iso8601Time
}

func (TapdIterationTask) TableName() string {
	return "_tool_tapd_iteration_tasks"
}
