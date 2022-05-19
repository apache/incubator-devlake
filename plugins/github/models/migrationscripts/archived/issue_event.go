package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GithubIssueEvent struct {
	GithubId        int       `gorm:"primaryKey"`
	IssueId         int       `gorm:"index;comment:References the Issue"`
	Type            string    `gorm:"type:varchar(255);comment:Events that can occur to an issue, ex. assigned, closed, labeled, etc."`
	AuthorUsername  string    `gorm:"type:varchar(255)"`
	GithubCreatedAt time.Time `gorm:"index"`
	archived.NoPKModel
}

func (GithubIssueEvent) TableName() string {
	return "_tool_github_issue_events"
}
