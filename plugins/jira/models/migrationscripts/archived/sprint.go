package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type JiraSprint struct {
	SourceId      uint64 `gorm:"primaryKey"`
	SprintId      uint64 `gorm:"primaryKey"`
	Self          string `gorm:"type:varchar(255)"`
	State         string `gorm:"type:varchar(255)"`
	Name          string `gorm:"type:varchar(255)"`
	StartDate     *time.Time
	EndDate       *time.Time
	CompleteDate  *time.Time
	OriginBoardID uint64
	archived.NoPKModel
}

type JiraBoardSprint struct {
	archived.NoPKModel
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	SprintId uint64 `gorm:"primaryKey"`
}

type JiraSprintIssue struct {
	archived.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	SprintId         uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *time.Time
	IssueCreatedDate *time.Time
}

func (JiraSprint) TableName() string {
	return "_tool_jira_sprints"
}

func (JiraBoardSprint) TableName() string {
	return "_tool_jira_board_sprints"
}

func (JiraSprintIssue) TableName() string {
	return "_tool_jira_sprint_issues"
}
