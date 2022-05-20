package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
)

type TapdIterationBug struct {
	common.NoPKModel
	ConnectionId   uint64 `gorm:"primaryKey"`
	IterationId    uint64 `gorm:"primaryKey"`
	WorkspaceID    uint64 `gorm:"primaryKey"`
	BugId          uint64 `gorm:"primaryKey"`
	ResolutionDate *helper.CSTTime
	BugCreatedDate *helper.CSTTime
}

func (TapdIterationBug) TableName() string {
	return "_tool_tapd_iteration_bugs"
}
