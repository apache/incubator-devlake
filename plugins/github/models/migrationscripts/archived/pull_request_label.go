package archived

import "github.com/merico-dev/lake/models/migrationscripts/archived"

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubPullRequestLabel struct {
	PullId    int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (GithubPullRequestLabel) TableName() string {
	return "_tool_github_pull_request_labels"
}
