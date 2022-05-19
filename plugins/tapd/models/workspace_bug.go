package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TapdWorkSpaceBug struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	BugId        uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkSpaceBug) TableName() string {
	return "_tool_tapd_workspace_bugs"
}
