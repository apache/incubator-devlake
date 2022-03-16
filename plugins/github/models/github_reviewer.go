package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
)

type GithubReviewer struct {
	GithubId      int `gorm:"primaryKey"`
	Login         string
	PullRequestId int `gorm:"primaryKey"`

	common.NoPKModel
	helper.RawDataOrigin
}
