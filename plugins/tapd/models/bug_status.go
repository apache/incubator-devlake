package models

import "github.com/merico-dev/lake/models/common"

type TapdBugStatus struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	EnglishName  string `gorm:"primaryKey;type:varchar(255)"`
	ChineseName  string
	IsLastStep   bool
	common.NoPKModel
}

func (TapdBugStatus) TableName() string {
	return "_tool_tapd_bug_statuses"
}
