package archived

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
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
