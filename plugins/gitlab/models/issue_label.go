package models

import (
"github.com/apache/incubator-devlake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GitlabIssueLabel struct {
	IssueId   int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (GitlabIssueLabel) TableName() string{
	return "_tool_gitlab_issue_labels"
}
