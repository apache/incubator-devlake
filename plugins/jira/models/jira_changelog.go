package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type JiraChangelog struct {
	common.NoPKModel

	// collected fields
	SourceId          uint64 `gorm:"primaryKey"`
	ChangelogId       uint64 `gorm:"primarykey"`
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	AuthorActive      bool
	Created           time.Time `gorm:"index"`
}

type JiraChangelogItem struct {
	common.NoPKModel

	// collected fields
	SourceId    uint64 `gorm:"primaryKey"`
	ChangelogId uint64 `gorm:"primaryKey"`
	Field       string `gorm:"primaryKey"`
	FieldType   string
	FieldId     string
	From        string
	FromString  string
	To          string
	ToString    string
}
