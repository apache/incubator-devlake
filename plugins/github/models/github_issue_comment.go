package models

import (
	"github.com/merico-dev/lake/plugins/helper"
	"time"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index;comment:References the Issue"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	helper.RawDataOrigin
}
