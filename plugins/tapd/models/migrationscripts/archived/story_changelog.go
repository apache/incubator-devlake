package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdStoryChangelog struct {
	ConnectionId   uint64            `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             uint64            `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id,string"`
	WorkspaceID    uint64            `json:"workspace_id,string"`
	WorkitemTypeID uint64            `json:"workitem_type_id,string"`
	Creator        string            `json:"creator"`
	Created        *core.Iso8601Time `json:"created"`
	ChangeSummary  string            `json:"change_summary"`
	Comment        string            `json:"comment"`
	EntityType     string            `json:"entity_type"`
	ChangeType     string            `json:"change_type"`
	StoryID        uint64            `json:"story_id,string"`
	common.NoPKModel
	FieldChanges []TapdStoryChangelogItem `json:"field_changes" gorm:"-"`
}

type TapdStoryChangelogItem struct {
	ConnectionId      uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdStoryChangelog) TableName() string {
	return "_tool_tapd_story_changelogs"
}
func (TapdStoryChangelogItem) TableName() string {
	return "_tool_tapd_story_changelog_items"
}
