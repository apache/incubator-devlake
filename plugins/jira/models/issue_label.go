package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type JiraIssueLabel struct {
	ConnectionId uint64 `gorm:"primaryKey;autoIncrement:false"`
	IssueId      uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName    string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JiraIssueLabel) TableName() string {
	return "_tool_jira_issue_labels"
}
