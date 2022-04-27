package models

import "github.com/merico-dev/lake/models/common"

type TapdStoryStatus struct {
	SourceId    Uint64s `gorm:"primaryKey"`
	WorkspaceID Uint64s `gorm:"primaryKey"`
	EnglishName string  `gorm:"primaryKey"`
	ChineseName string
	IsLastStep  bool
	common.NoPKModel
}

func (TapdStoryStatus) TableName() string {
	return "_tool_tapd_story_statuses"
}
