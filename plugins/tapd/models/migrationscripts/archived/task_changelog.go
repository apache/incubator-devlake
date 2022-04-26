package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdTaskChangelog struct {
	SourceId       models.Uint64s    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             models.Uint64s    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId    models.Uint64s    `json:"workspace_id"`
	WorkitemTypeID models.Uint64s    `json:"workitem_type_id"`
	Creator        string            `json:"creator"`
	Created        *core.Iso8601Time `json:"created"`
	ChangeSummary  string            `json:"change_summary"`
	Comment        string            `json:"comment"`
	EntityType     string            `json:"entity_type"`
	ChangeType     string            `json:"change_type"`
	ChangeTypeText string            `json:"change_type_text"`
	TaskID         models.Uint64s    `json:"task_id"`
	common.NoPKModel
	FieldChanges []TapdTaskChangelogItem `json:"field_changes" gorm:"-"`
}

type TapdTaskChangelogItem struct {
	SourceId          models.Uint64s `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       models.Uint64s `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string         `json:"field" gorm:"primaryKey"`
	ValueBeforeParsed string         `json:"value_before"`
	ValueAfterParsed  string         `json:"value_after"`
	IterationIdFrom   models.Uint64s
	IterationIdTo     models.Uint64s
	common.NoPKModel
}

func (TapdTaskChangelog) TableName() string {
	return "_tool_tapd_task_changelogs"
}
func (TapdTaskChangelogItem) TableName() string {
	return "_tool_tapd_task_changelog_items"
}
