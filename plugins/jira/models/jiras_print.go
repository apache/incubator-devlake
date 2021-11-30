package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type JiraSprint struct {
	common.NoPKModel
	SourceId      uint64 `gorm:"primaryKey"`
	SprintId      uint64 `gorm:"primaryKey"`
	Self          string
	State         string
	Name          string
	StartDate     *time.Time
	EndDate       *time.Time
	CompleteDate  *time.Time
	OriginBoardID uint64
}

type JiraBoardSprint struct {
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	SprintId uint64 `gorm:"primaryKey"`
}

type JiraSprintIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	SprintId uint64 `gorm:"primaryKey"`
	IssueId  uint64 `gorm:"primaryKey"`
}
