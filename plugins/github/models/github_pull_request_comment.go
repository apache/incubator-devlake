package models

import (
	"github.com/merico-dev/lake/plugins/helper"
	"time"
)

type GithubPullRequestComment struct {
	GithubId        int `gorm:"primaryKey"`
	PullRequestId   int `gorm:"index"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time `gorm:"index"`

	helper.RawDataOrigin
}
