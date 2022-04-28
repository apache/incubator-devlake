package models

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdTaskChangelog struct {
	SourceId       Uint64s           `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             Uint64s           `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceID    Uint64s           `json:"workspace_id"`
	WorkitemTypeID Uint64s           `json:"workitem_type_id"`
	Creator        string            `json:"creator"`
	Created        *core.Iso8601Time `json:"created"`
	ChangeSummary  string            `json:"change_summary"`
	Comment        string            `json:"comment"`
	EntityType     string            `json:"entity_type"`
	ChangeType     string            `json:"change_type"`
	ChangeTypeText string            `json:"change_type_text"`
	TaskID         Uint64s           `json:"task_id"`
	common.NoPKModel
	FieldChanges []TapdTaskChangelogItemRes `json:"field_changes" gorm:"-"`
}

type TapdTaskChangelogItem struct {
	SourceId          Uint64s `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       Uint64s `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string  `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed string  `json:"value_before_parsed"`
	ValueAfterParsed  string  `json:"value_after_parsed"`
	IterationIdFrom   Uint64s
	IterationIdTo     Uint64s
	common.NoPKModel
}

type TapdTaskChangelogItemRes struct {
	SourceId          Uint64s         `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       Uint64s         `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string          `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed json.RawMessage `json:"value_before_parsed"`
	ValueAfterParsed  json.RawMessage `json:"value_after_parsed"`
	IterationIdFrom   Uint64s
	IterationIdTo     Uint64s
	common.NoPKModel
}

func (TapdTaskChangelog) TableName() string {
	return "_tool_tapd_task_changelogs"
}
func (TapdTaskChangelogItem) TableName() string {
	return "_tool_tapd_task_changelog_items"
}
