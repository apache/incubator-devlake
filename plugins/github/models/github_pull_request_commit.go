package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubPullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey;type:varchar(40)"`
	PullRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	common.NoPKModel
}

func (GithubPullRequestCommit) TableName() string{
	return "_tool_github_pull_request_commits"
}

