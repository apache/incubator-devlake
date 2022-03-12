package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubPullRequestComment struct {
	GithubId        int `gorm:"primaryKey"`
	PullRequestId   int `gorm:"index"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time

	common.NoPKModel
}
