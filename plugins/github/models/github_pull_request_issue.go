package models

type GithubPullRequestIssue struct {
	PullRequestId int `gorm:"primaryKey"`
	IssueId       int `gorm:"primaryKey"`
	PullNumber    int
	IssueNumber   int
}
