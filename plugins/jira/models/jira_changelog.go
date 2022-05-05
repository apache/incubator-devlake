package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type JiraChangelog struct {
	common.NoPKModel

	// collected fields
	ConnectionId      uint64 `gorm:"primaryKey"`
	ChangelogId       uint64 `gorm:"primarykey"`
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string `gorm:"type:varchar(255)"`
	AuthorDisplayName string `gorm:"type:varchar(255)"`
	AuthorActive      bool
	Created           time.Time `gorm:"index"`
}

type JiraChangelogItem struct {
	common.NoPKModel

	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	ChangelogId  uint64 `gorm:"primaryKey"`
	Field        string `gorm:"primaryKey"`
	FieldType    string
	FieldId      string
	From         string
	FromString   string
	To           string
	ToString     string
}

func (JiraChangelog) TableName() string {
	return "_tool_jira_changelogs"
}

func (JiraChangelogItem) TableName() string {
	return "_tool_jira_changelog_items"
}
