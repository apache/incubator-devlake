package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index;comment:References the Pull Request"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time

	models.NoPKModel
}
