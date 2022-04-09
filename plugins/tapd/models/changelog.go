package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdChangelog struct {
	SourceId       uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId    uint64     `json:"workspace_id"`
	WorkitemTypeID uint64     `json:"workitem_type_id"`
	Creator        string     `json:"creator"`
	Created        *time.Time `json:"created"`
	ChangeSummary  string     `json:"change_summary"`
	Comment        string     `json:"comment"`
	EntityType     string     `json:"entity_type"`
	ChangeType     string     `json:"change_type"`
	ChangeTypeText string     `json:"change_type_text"`
	IssueId        uint64     `gorm:"index"`
	common.NoPKModel
}

type TapdChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed string `json:"value_before"`
	ValueAfterParsed  string `json:"value_after"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

type TapdBugChangelogApiRes struct {
	ID          string `json:"id"`
	WorkspaceId string `json:"workspace_id"`
	BugID       string `json:"bug_id"`
	Author      string `json:"author"`
	Field       string `json:"field"`
	OldValue    string `json:"old_value"`
	NewValue    string `json:"new_value"`
	Memo        string `json:"memo"`
	Created     string `json:"created"`
}
type TapdBugChangelog struct {
	ID          uint64     `json:"id"`
	WorkspaceId uint64     `json:"workspace_id"`
	BugID       uint64     `json:"bug_id"`
	Author      string     `json:"author"`
	Field       string     `json:"field"`
	OldValue    string     `json:"old_value"`
	NewValue    string     `json:"new_value"`
	Memo        string     `json:"memo"`
	Created     *time.Time `json:"created"`
}

type TapdTaskChangelogApiRes struct {
	ID             string                `json:"id"`
	WorkspaceId    string                `json:"workspace_id"`
	WorkitemTypeID string                `json:"workitem_type_id"`
	Creator        string                `json:"creator"`
	Created        string                `json:"created"`
	ChangeSummary  string                `json:"change_summary"`
	Comment        string                `json:"comment"`
	FieldChanges   []ChangelogItemApiRes `json:"field_changes"`
	EntityType     string                `json:"entity_type"`
	ChangeType     string                `json:"change_type"`
	ChangeTypeText string                `json:"change_type_text"`
	TaskID         string                `json:"task_id"`
}
type TapdTaskChangelog struct {
	ID             uint64                `json:"id"`
	WorkspaceId    uint64                `json:"workspace_id"`
	WorkitemTypeID uint64                `json:"workitem_type_id"`
	Creator        string                `json:"creator"`
	Created        *time.Time            `json:"created"`
	ChangeSummary  string                `json:"change_summary"`
	Comment        string                `json:"comment"`
	FieldChanges   []ChangelogItemApiRes `json:"field_changes"`
	EntityType     string                `json:"entity_type"`
	ChangeType     string                `json:"change_type"`
	ChangeTypeText string                `json:"change_type_text"`
	TaskID         uint64                `json:"task_id"`
}
type TapdStoryChangelogApiRes struct {
	ID             string                `json:"id"`
	WorkspaceId    string                `json:"workspace_id"`
	WorkitemTypeID string                `json:"workitem_type_id"`
	Creator        string                `json:"creator"`
	Created        string                `json:"created"`
	ChangeSummary  string                `json:"change_summary"`
	Comment        string                `json:"comment"`
	FieldChanges   []ChangelogItemApiRes `json:"field_changes"`
	EntityType     string                `json:"entity_type"`
	StoryID        string                `json:"story_id"`
}
type TapdStoryChangelog struct {
	ID             uint64                `json:"id"`
	WorkspaceId    uint64                `json:"workspace_id"`
	WorkitemTypeID uint64                `json:"workitem_type_id"`
	Creator        string                `json:"creator"`
	Created        *time.Time            `json:"created"`
	ChangeSummary  string                `json:"change_summary"`
	Comment        string                `json:"comment"`
	FieldChanges   []ChangelogItemApiRes `json:"field_changes"`
	EntityType     string                `json:"entity_type"`
	StoryID        uint64                `json:"story_id"`
}

type ChangelogItemApiRes struct {
	Field             string `json:"field"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	FieldLabel        string `json:"field_label"`
}

type ChangelogTmp struct {
	Id              uint64
	IssueId         uint64
	AuthorId        string
	AuthorName      string
	FieldId         string
	FieldName       string
	From            string
	To              string
	IterationIdFrom uint64
	IterationIdTo   uint64
	CreatedDate     time.Time
	common.RawDataOrigin
}

func (TapdChangelog) TableName() string {
	return "_tool_tapd_changelogs"
}
func (TapdChangelogItem) TableName() string {
	return "_tool_tapd_changelog_items"
}
