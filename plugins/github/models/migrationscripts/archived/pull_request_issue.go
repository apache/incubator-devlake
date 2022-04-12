package archived

import "github.com/merico-dev/lake/models/common"

type GithubPullRequestIssue struct {
	PullRequestId int `gorm:"primaryKey"`
	IssueId       int `gorm:"primaryKey"`
	PullNumber    int
	IssueNumber   int
	common.NoPKModel
}

func (GithubPullRequestIssue) TableName() string{
	return "_tool_github_pull_request_issues"
}

