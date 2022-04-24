package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdTaskChangelog struct {
	SourceId       uint64    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             uint64    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId    uint64    `json:"workspace_id"`
	WorkitemTypeID uint64    `json:"workitem_type_id"`
	Creator        string    `json:"creator"`
	Created        time.Time `json:"created"`
	ChangeSummary  string    `json:"change_summary"`
	Comment        string    `json:"comment"`
	EntityType     string    `json:"entity_type"`
	ChangeType     string    `json:"change_type"`
	ChangeTypeText string    `json:"change_type_text"`
	TaskID         uint64    `json:"task_id"`
	common.NoPKModel
}

type TapdTaskChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed string `json:"value_before"`
	ValueAfterParsed  string `json:"value_after"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

type TapdTaskChangelogApiRes struct {
	ID             string                    `json:"id"`
	WorkspaceId    string                    `json:"workspace_id"`
	WorkitemTypeID string                    `json:"workitem_type_id"`
	Creator        string                    `json:"creator"`
	Created        string                    `json:"created"`
	ChangeSummary  string                    `json:"change_summary"`
	Comment        string                    `json:"comment"`
	FieldChanges   []TaskChangelogItemApiRes `json:"field_changes"`
	EntityType     string                    `json:"entity_type"`
	ChangeType     string                    `json:"change_type"`
	ChangeTypeText string                    `json:"change_type_text"`
	TaskID         string                    `json:"task_id"`
}
type TaskChangelogItemApiRes struct {
	Field             string `json:"field"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	FieldLabel        string `json:"field_label"`
}

func (TapdTaskChangelog) TableName() string {
	return "_tool_tapd_task_changelogs"
}
func (TapdTaskChangelogItem) TableName() string {
	return "_tool_tapd_task_changelog_items"
}
