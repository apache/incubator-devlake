package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
)

type GithubPullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey"`
	PullRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	common.NoPKModel
	helper.RawDataOrigin
}
