package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GithubPullRequestCommit struct {
	Sha            string `gorm:"primaryKey"`
	PullRequestId  int    `gorm:"index"` // This value links to pull request
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   sql.NullTime
	CommitterName  string
	CommitterEmail string
	CommittedDate  sql.NullTime
	Message        string
	Url            string

	models.NoPKModel
}
