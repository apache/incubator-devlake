package models

import "github.com/merico-dev/lake/models/common"

type TapdStoryStatus struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	EnglishName  string `gorm:"primaryKey;type:varchar(255)"`
	ChineseName  string
	IsLastStep   bool
	common.NoPKModel
}

func (TapdStoryStatus) TableName() string {
	return "_tool_tapd_story_statuses"
}
