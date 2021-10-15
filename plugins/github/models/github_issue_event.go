package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GithubIssueEvent struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index"` // This value links to pull request
	Type            string
	AuthorUsername  string
	GithubCreatedAt sql.NullTime

	models.NoPKModel
}
