package archived

import (
	"time"
)

type Sprint struct {
	DomainEntity
	Name            string `gorm:"type:char(255)"`
	Url             string `gorm:"type:char(255)"`
	Status          string `gorm:"type:char(100)"`
	StartedDate     *time.Time
	EndedDate       *time.Time
	CompletedDate   *time.Time
	OriginalBoardID string `gorm:"type:char(255)"`
}

type SprintIssue struct {
	NoPKModel
	SprintId      string `gorm:"primaryKey;type:varchar(255)"`
	IssueId       string `gorm:"primaryKey;type:varchar(255)"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    *string `gorm:"type:varchar(255)"`
	ResolvedStage *string `gorm:"type:varchar(255)"`
}
