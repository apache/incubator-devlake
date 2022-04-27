package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdStoryStatus struct {
	SourceId    models.Uint64s `gorm:"primaryKey"`
	WorkspaceID models.Uint64s `gorm:"primaryKey"`
	EnglishName string         `gorm:"primaryKey"`
	ChineseName string
	IsLastStep  bool
	common.NoPKModel
}

func (TapdStoryStatus) TableName() string {
	return "_tool_tapd_story_statuses"
}
