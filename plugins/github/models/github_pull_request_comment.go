package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubPullRequestComment struct {
	GithubId        int `gorm:"primaryKey"`
	PullRequestId   int `gorm:"index"` // This value links to pull request
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time

	models.NoPKModel
}
