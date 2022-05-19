package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type GithubIssueEvent struct {
	GithubId        int       `gorm:"primaryKey"`
	IssueId         int       `gorm:"index;comment:References the Issue"`
	Type            string    `gorm:"type:varchar(255);comment:Events that can occur to an issue, ex. assigned, closed, labeled, etc."`
	AuthorUsername  string    `gorm:"type:varchar(255)"`
	GithubCreatedAt time.Time `gorm:"index"`
	common.NoPKModel
}

func (GithubIssueEvent) TableName() string {
	return "_tool_github_issue_events"
}
