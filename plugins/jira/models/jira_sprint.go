package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type JiraSprint struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	SprintId      uint64 `gorm:"primaryKey"`
	Self          string `gorm:"type:varchar(255)"`
	State         string `gorm:"type:varchar(255)"`
	Name          string `gorm:"type:varchar(255)"`
	StartDate     *time.Time
	EndDate       *time.Time
	CompleteDate  *time.Time
	OriginBoardID uint64
	common.NoPKModel
}

type JiraBoardSprint struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`
	SprintId     uint64 `gorm:"primaryKey"`
}

type JiraSprintIssue struct {
	common.NoPKModel
	ConnectionId     uint64 `gorm:"primaryKey"`
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
