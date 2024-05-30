package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// IssueStatusHistory records issue status history (status original value)
// end_date of current status is set to now() to avoid false assumption of future status.
// handled by ConvertIssueStatusHistory task
type IssueStatusHistory struct {
	common.NoPKModel
	IssueId           string     `gorm:"primaryKey;type:varchar(255)"`
	Status            string     `gorm:"type:varchar(100)"`
	OriginalStatus    string     `gorm:"primaryKey;type:varchar(255)"`
	StartDate         time.Time  `gorm:"primaryKey"`
	EndDate           *time.Time `gorm:"type:timestamp"`
	IsCurrentStatus   bool       `gorm:"type:boolean"`
	IsFirstStatus     bool       `gorm:"type:boolean"`
	StatusTimeMinutes int32      `gorm:"type:integer"`
}

func (IssueStatusHistory) TableName() string {
	return "issue_status_history"
}
