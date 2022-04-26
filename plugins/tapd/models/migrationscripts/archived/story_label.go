package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdStoryLabel struct {
	StoryId   models.Uint64s `gorm:"primaryKey;autoIncrement:false"`
	LabelName string         `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (TapdStoryLabel) TableName() string {
	return "_tool_tapd_story_labels"
}
