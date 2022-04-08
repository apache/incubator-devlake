package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubReviewer struct {
	GithubId      int    `gorm:"primaryKey"`
	Login         string `gorm:"type:varchar(255)"`
	PullRequestId int    `gorm:"primaryKey"`

	common.NoPKModel
}
