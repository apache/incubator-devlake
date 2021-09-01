package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type JiraChangelog struct {
	models.Model

	// collected fields
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	AuthorActive      bool
	Created           time.Time
}

type JiraChangelogItem struct {
	models.Model

	ChangelogId uint64 `gorm:"index"`
	Field       string
	FieldType   string
	FieldId     string
	From        string
	FromString  string
	To          string
	ToString    string
}
