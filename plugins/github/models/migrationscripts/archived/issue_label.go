package archived

import "github.com/merico-dev/lake/models/migrationscripts/archived"

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubIssueLabel struct {
	IssueId   int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (GithubIssueLabel) TableName() string {
	return "_tool_github_issue_labels"
}
