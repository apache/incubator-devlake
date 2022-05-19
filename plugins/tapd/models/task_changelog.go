package models

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdTaskChangelog struct {
	ConnectionId   uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID             uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	WorkspaceID    uint64          `json:"workspace_id,string"`
	WorkitemTypeID uint64          `json:"workitem_type_id,string"`
	Creator        string          `json:"creator"`
	Created        *helper.CSTTime `json:"created"`
	ChangeSummary  string          `json:"change_summary"`
	Comment        string          `json:"comment"`
	EntityType     string          `json:"entity_type"`
	ChangeType     string          `json:"change_type"`
	ChangeTypeText string          `json:"change_type_text"`
	TaskID         uint64          `json:"task_id,string"`
	common.NoPKModel
	FieldChanges []TapdTaskChangelogItemRes `json:"field_changes" gorm:"-"`
}

type TapdTaskChangelogItem struct {
	ConnectionId      uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

type TapdTaskChangelogItemRes struct {
	ConnectionId      uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ChangelogId       uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string          `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed json.RawMessage `json:"value_before_parsed"`
	ValueAfterParsed  json.RawMessage `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdTaskChangelog) TableName() string {
	return "_tool_tapd_task_changelogs"
}
func (TapdTaskChangelogItem) TableName() string {
	return "_tool_tapd_task_changelog_items"
}
