package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GithubPullRequestComment struct {
	GithubId        int `gorm:"primaryKey"`
	PullRequestId   int `gorm:"index"`
	Body            string
	AuthorUsername  string `gorm:"type:varchar(255)"`
	AuthorUserId    int
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	archived.NoPKModel
}

func (GithubPullRequestComment) TableName() string {
	return "_tool_github_pull_request_comments"
}
