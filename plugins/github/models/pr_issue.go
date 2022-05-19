package models

import "github.com/apache/incubator-devlake/models/common"

type GithubPullRequestIssue struct {
	PullRequestId     int `gorm:"primaryKey"`
	IssueId           int `gorm:"primaryKey"`
	PullRequestNumber int
	IssueNumber       int
	common.NoPKModel
}

func (GithubPullRequestIssue) TableName() string {
	return "_tool_github_pull_request_issues"
}
