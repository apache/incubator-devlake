package archived

import (
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GitlabIssueLabel struct {
	IssueId   int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (GitlabIssueLabel) TableName() string{
	return "_tool_gitlab_issue_labels"
}
