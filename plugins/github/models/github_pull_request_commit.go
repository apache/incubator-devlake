package models

import "github.com/merico-dev/lake/models"

type GithubPullRequestCommit struct {
	Sha            string `gorm:"primaryKey"`
	PullRequestId  int    `gorm:"index"` // This value links to pull request
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   string
	CommitterName  string
	CommitterEmail string
	CommittedDate  string
	Message        string
	Url            string

	models.NoPKModel
}
