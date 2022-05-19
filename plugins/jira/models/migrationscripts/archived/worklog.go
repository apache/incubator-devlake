package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type JiraWorklog struct {
	archived.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primarykey"`
	WorklogId        string `gorm:"primarykey;type:varchar(255)"`
	AuthorId         string `gorm:"type:varchar(255)"`
	UpdateAuthorId   string `gorm:"type:varchar(255)"`
	TimeSpent        string `gorm:"type:varchar(255)"`
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}

func (JiraWorklog) TableName() string {
	return "_tool_jira_worklogs"
}
