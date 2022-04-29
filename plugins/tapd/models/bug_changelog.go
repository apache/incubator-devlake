package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdBugChangelog struct {
	SourceId    uint64            `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	WorkspaceID uint64            `gorm:"type:INT(10) UNSIGNED NOT NULL"`
	ID          uint64            `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id,string"`
	BugID       uint64            `json:"bug_id,string"`
	Author      string            `json:"author"`
	Field       string            `gorm:"primaryKey;type:varchar(255)" json:"field"`
	OldValue    string            `json:"old_value"`
	NewValue    string            `json:"new_value"`
	Memo        string            `json:"memo"`
	Created     *core.Iso8601Time `json:"created"`
	common.NoPKModel
}

type TapdBugChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdBugChangelog) TableName() string {
	return "_tool_tapd_bug_changelogs"
}
func (TapdBugChangelogItem) TableName() string {
	return "_tool_tapd_bug_changelog_items"
}
