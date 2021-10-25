package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubPullRequestCommit struct {
	Sha            string `gorm:"primaryKey"`
	PullRequestId  int    `gorm:"index"` // This value links to pull request
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	Message        string
	Url            string

	models.NoPKModel
}
