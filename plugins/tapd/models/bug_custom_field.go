package models

import "github.com/apache/incubator-devlake/models/common"

type TapdBugCustomFields struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	WorkspaceID  uint64 `json:"workspace_id,string"`
	EntryType    string `json:"entry_type" gorm:"type:varchar(20)"`
	CustomField  string `json:"custom_field" gorm:"type:varchar(255)"`
	Type         string `json:"type" gorm:"type:varchar(20)"`
	Name         string `json:"name" gorm:"type:varchar(255)"`
	Options      string `json:"options" gorm:"type:text"`
	Enabled      string `json:"enabled" gorm:"type:varchar(255)"`
	Sort         string `json:"sort" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (TapdBugCustomFields) TableName() string {
	return "_tool_tapd_bug_custom_fields"
}
