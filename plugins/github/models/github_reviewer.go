package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubReviewer struct {
	GithubId      int `gorm:"primaryKey"`
	Login         string
	PullRequestId int

	common.NoPKModel
}
