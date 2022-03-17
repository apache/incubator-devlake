package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubPullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey"`
	PullRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	common.NoPKModel
}
