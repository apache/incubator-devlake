package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdIterationBug struct {
	common.NoPKModel
	SourceId       uint64 `gorm:"primaryKey"`
	IterationId    uint64 `gorm:"primaryKey"`
	WorkspaceID    uint64 `gorm:"primaryKey"`
	BugId          uint64 `gorm:"primaryKey"`
	ResolutionDate *core.Iso8601Time
	BugCreatedDate *core.Iso8601Time
}

func (TapdIterationBug) TableName() string {
	return "_tool_tapd_iteration_bugs"
}
