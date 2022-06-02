package archived

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdStoryCategory struct {
	ConnectionId uint64         `gorm:"primaryKey"`
	ID           uint64         `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name         string         `json:"name" gorm:"type:varchar(255)"`
	Description  string         `json:"description"`
	ParentID     uint64         `json:"parent_id,string"`
	Created      helper.CSTTime `json:"created"`
	Modified     helper.CSTTime `json:"modified"`
	common.NoPKModel
}

func (TapdStoryCategory) TableName() string {
	return "_tool_tapd_story_categories"
}
