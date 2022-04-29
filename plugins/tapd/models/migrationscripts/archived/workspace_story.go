package archived

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdWorkSpaceStory struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	StoryId      uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkSpaceStory) TableName() string {
	return "_tool_tapd_workspace_stories"
}
