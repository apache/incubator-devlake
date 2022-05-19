package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubPullRequestLabel struct {
	PullId    int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (GithubPullRequestLabel) TableName() string {
	return "_tool_github_pull_request_labels"
}
