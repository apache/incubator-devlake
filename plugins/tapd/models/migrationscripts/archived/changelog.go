package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdChangelog struct {
	SourceId       uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID             uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId    uint64     `json:"workspace_id"`
	WorkitemTypeID uint64     `json:"workitem_type_id"`
	Creator        string     `gorm:"type:varchar(255)"`
	Created        *time.Time `json:"created"`
	ChangeSummary  string     `gorm:"type:varchar(255)"`
	Comment        string     `gorm:"type:varchar(255)"`
	EntityType     string     `gorm:"type:varchar(255)"`
	ChangeType     string     `gorm:"type:varchar(255)"`
	ChangeTypeText string     `gorm:"type:varchar(255)"`
	IssueId        uint64     `gorm:"index"`
	common.NoPKModel
}

type TapdChangelogItem struct {
	SourceId          uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	Field             string `gorm:"type:varchar(255);primaryKey"`
	ValueBeforeParsed string `gorm:"type:varchar(255)"`
	ValueAfterParsed  string `gorm:"type:varchar(255)"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdChangelog) TableName() string {
	return "_tool_tapd_changelogs"
}
func (TapdChangelogItem) TableName() string {
	return "_tool_tapd_changelog_items"
}
