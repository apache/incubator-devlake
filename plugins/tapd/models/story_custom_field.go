package models

import "github.com/apache/incubator-devlake/models/common"

type TapdStoryCustomFields struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	WorkspaceID  uint64 `json:"workspace_id,string"`
	EntryType    string `json:"entry_type"`
	CustomField  string `json:"custom_field"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Options      string `json:"options"`
	Enabled      string `json:"enabled"`
	Sort         string `json:"sort"`
	common.NoPKModel
}

func (TapdStoryCustomFields) TableName() string {
	return "_tool_tapd_story_custom_fields"
}
