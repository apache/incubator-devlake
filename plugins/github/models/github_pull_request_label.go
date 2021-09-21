package models

import "github.com/merico-dev/lake/models"

type GithubPullRequestLabel struct {
	GithubId      int `gorm:"primaryKey"`
	PullRequestId int `gorm:"index"`
	Name          string
	Description   string
	Url           string

	models.NoPKModel
}
