package models

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdTaskChangelog struct {
	SourceId       uint64            `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             uint64            `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id,string"`
	WorkspaceID    uint64            `json:"workspace_id,string"`
	WorkitemTypeID uint64            `json:"workitem_type_id,string"`
	Creator        string            `json:"creator"`
	Created        *core.Iso8601Time `json:"created"`
	ChangeSummary  string            `json:"change_summary"`
	Comment        string            `json:"comment"`
	EntityType     string            `json:"entity_type"`
	ChangeType     string            `json:"change_type"`
	ChangeTypeText string            `json:"change_type_text"`
	TaskID         uint64            `json:"task_id,string"`
	common.NoPKModel
	FieldChanges []TapdTaskChangelogItemRes `json:"field_changes" gorm:"-"`
}

type TapdTaskChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

type TapdTaskChangelogItemRes struct {
	SourceId          uint64          `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64          `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
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
