package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type TapdIssueStatusHistory struct {
	common.NoPKModel
	SourceId       uint64 `gorm:"primaryKey"`
	WorkspaceId    uint64
	IssueId        uint64    `gorm:"primaryKey"`
	OriginalStatus string    `gorm:"primaryKey"`
	StartDate      time.Time `gorm:"primaryKey"`
	EndDate        time.Time
}

type TapdIssueAssigneeHistory struct {
	common.NoPKModel
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceId uint64

	IssueId   uint64    `gorm:"primaryKey"`
	Assignee  string    `gorm:"primaryKey"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   time.Time
}

type TapdIssueSprintsHistory struct {
	common.NoPKModel
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceId uint64

	IssueId   uint64    `gorm:"primaryKey"`
	SprintId  uint64    `gorm:"primaryKey"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   time.Time
}
