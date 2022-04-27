package models

import "github.com/merico-dev/lake/models/common"

type TapdBugStatus struct {
	SourceId    Uint64s `gorm:"primaryKey"`
	WorkspaceID Uint64s `gorm:"primaryKey"`
	EnglishName string  `gorm:"primaryKey"`
	ChineseName string
	IsLastStep  bool
	common.NoPKModel
}

func (TapdBugStatus) TableName() string {
	return "_tool_tapd_bug_statuses"
}
