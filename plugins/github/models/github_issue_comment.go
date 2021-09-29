package models

import "github.com/merico-dev/lake/models"

type GithubIssueComment struct {
	GithubId       int `gorm:"primaryKey"`
	IssueId        int `gorm:"index"` // This value links to pull request
	Body           string
	AuthorUsername string

	models.NoPKModel
}
