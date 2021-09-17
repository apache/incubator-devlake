package models

import "github.com/merico-dev/lake/models"

type GithubPullRequest struct {
	GithubId        int `gorm:"primaryKey"`
	RepositoryId    int `gorm:"index"`
	State           string
	Title           string
	HTMLUrl         string
	MergedAt        string
	GithubCreatedAt string
	ClosedAt        string

	models.NoPKModel
}
