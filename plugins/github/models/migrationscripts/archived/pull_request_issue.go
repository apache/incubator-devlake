package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type GithubPullRequestIssue struct {
	PullRequestId     int `gorm:"primaryKey"`
	IssueId           int `gorm:"primaryKey"`
	PullRequestNumber int
	IssueNumber       int
	archived.NoPKModel
}

func (GithubPullRequestIssue) TableName() string {
	return "_tool_github_pull_request_issues"
}
