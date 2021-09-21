package models

import "github.com/merico-dev/lake/models"

type GithubReviewer struct {
	GithubId      int `gorm:"primaryKey"`
	Login         string
	PullRequestId int

	models.NoPKModel
}
