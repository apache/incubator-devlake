package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdBugChangelog struct {
	SourceId    Uint64s           `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	WorkspaceId Uint64s           `gorm:"type:INT(10) UNSIGNED NOT NULL"`
	ID          Uint64s           `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	BugID       Uint64s           `json:"bug_id"`
	Author      string            `json:"author"`
	Field       string            `gorm:"primaryKey" json:"field"`
	OldValue    string            `json:"old_value"`
	NewValue    string            `json:"new_value"`
	Memo        string            `json:"memo"`
	Created     *core.Iso8601Time `json:"created"`
	common.NoPKModel
}

type TapdBugChangelogItem struct {
	SourceId          Uint64s `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       Uint64s `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string  `json:"field" gorm:"primaryKey;foreignKey:ChangelogId;references:ID"`
	ValueBeforeParsed string  `json:"value_before"`
	ValueAfterParsed  string  `json:"value_after"`
	IterationIdFrom   Uint64s
	IterationIdTo     Uint64s
	common.NoPKModel
}

func (TapdBugChangelog) TableName() string {
	return "_tool_tapd_bug_changelogs"
}
func (TapdBugChangelogItem) TableName() string {
	return "_tool_tapd_bug_changelog_items"
}
