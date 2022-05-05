package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdIterationStory struct {
	common.NoPKModel
	ConnectionId     uint64 `gorm:"primaryKey"`
	IterationId      uint64 `gorm:"primaryKey"`
	WorkspaceID      uint64 `gorm:"primaryKey"`
	StoryId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *core.CSTTime
	StoryCreatedDate *core.CSTTime
}

func (TapdIterationStory) TableName() string {
	return "_tool_tapd_iteration_stories"
}
