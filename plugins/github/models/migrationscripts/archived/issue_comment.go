package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index;comment:References the Issue"`
	Body            string
	AuthorUsername  string `gorm:"type:varchar(255)"`
	AuthorUserId    int
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	common.NoPKModel
}

func (GithubIssueComment) TableName() string {
	return "_tool_github_issue_comments"
}
