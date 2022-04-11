package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type IssueStatusHistory struct {
	common.NoPKModel
	IssueId        string    `gorm:"primaryKey;type:varchar(255)"`
	OriginalStatus string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate      time.Time `gorm:"primaryKey"`
	EndDate        *time.Time
}

func (IssueStatusHistory) TableName() string {
	return "issue_status_history"
}

type IssueAssigneeHistory struct {
	common.NoPKModel
	IssueId   string    `gorm:"primaryKey;type:varchar(255)"`
	Assignee  string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   *time.Time
}

func (IssueAssigneeHistory) TableName() string {
	return "issue_assignee_history"
}

type IssueSprintsHistory struct {
	common.NoPKModel
	IssueId   string    `gorm:"primaryKey;type:varchar(255)"`
	SprintId  string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   *time.Time
}

func (IssueSprintsHistory) TableName() string {
	return "issue_sprints_history"
}
