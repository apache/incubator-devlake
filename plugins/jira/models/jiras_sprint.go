package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type JiraSprint struct {
	SourceId      uint64 `gorm:"primaryKey"`
	SprintId      uint64 `gorm:"primaryKey"`
	Self          string
	State         string
	Name          string
	StartDate     *time.Time
	EndDate       *time.Time
	CompleteDate  *time.Time
	OriginBoardID uint64
	common.NoPKModel
}

type JiraBoardSprint struct {
	common.NoPKModel
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	SprintId uint64 `gorm:"primaryKey"`
}

type JiraSprintIssue struct {
	common.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	SprintId         uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *time.Time
	IssueCreatedDate *time.Time
}
