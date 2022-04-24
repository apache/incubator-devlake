package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdBugChangelog struct {
	SourceId    uint64    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	WorkspaceId uint64    `gorm:"type:INT(10) UNSIGNED NOT NULL"`
	ID          uint64    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	BugID       uint64    `json:"bug_id"`
	Author      string    `json:"author"`
	Field       string    `json:"field"`
	OldValue    string    `json:"old_value"`
	NewValue    string    `json:"new_value"`
	Memo        string    `json:"memo"`
	Created     time.Time `json:"created"`
	common.NoPKModel
}

type TapdBugChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed string `json:"value_before"`
	ValueAfterParsed  string `json:"value_after"`
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
