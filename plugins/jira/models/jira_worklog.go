package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type JiraWorklog struct {
	models.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primarykey"`
	WorklogId        string `gorm:"primarykey"`
	AuthorId         string
	UpdateAuthorId   string
	TimeSpent        string
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}
