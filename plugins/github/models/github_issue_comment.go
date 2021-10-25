package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index"` // This value links to pull request
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time

	models.NoPKModel
}
