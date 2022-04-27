package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type TapdIssueStatusHistory struct {
	common.NoPKModel
	SourceId       Uint64s `gorm:"primaryKey"`
	WorkspaceID    Uint64s
	IssueId        Uint64s   `gorm:"primaryKey"`
	OriginalStatus string    `gorm:"primaryKey;type:varchar(250)"`
	StartDate      time.Time `gorm:"primaryKey"`
	EndDate        time.Time
}

type TapdIssueAssigneeHistory struct {
	common.NoPKModel
	SourceId    Uint64s `gorm:"primaryKey"`
	WorkspaceID Uint64s

	IssueId   Uint64s   `gorm:"primaryKey"`
	Assignee  string    `gorm:"primaryKey;type:varchar(250)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   time.Time
}

type TapdIssueSprintHistory struct {
	common.NoPKModel
	SourceId    Uint64s `gorm:"primaryKey"`
	WorkspaceID Uint64s
	ChangelogId Uint64s
	IssueId     Uint64s   `gorm:"primaryKey"`
	SprintId    Uint64s   `gorm:"primaryKey"`
	StartDate   time.Time `gorm:"primaryKey"`
	EndDate     time.Time
}

func (TapdIssueStatusHistory) TableName() string {
	return "_tool_tapd_issue_status_histories"
}
func (TapdIssueAssigneeHistory) TableName() string {
	return "_tool_tapd_issue_assignee_histories"
}
func (TapdIssueSprintHistory) TableName() string {
	return "_tool_tapd_issue_sprint_histories"
}
