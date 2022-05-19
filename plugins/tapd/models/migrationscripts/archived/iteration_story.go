package archived

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdIterationStory struct {
	common.NoPKModel
	ConnectionId     uint64 `gorm:"primaryKey"`
	IterationId      uint64 `gorm:"primaryKey"`
	WorkspaceID      uint64 `gorm:"primaryKey"`
	StoryId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *helper.CSTTime
	StoryCreatedDate *helper.CSTTime
}

func (TapdIterationStory) TableName() string {
	return "_tool_tapd_iteration_stories"
}
