package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// IssueAssigneeHistory records issue assignee history
// end_date of current assignee is set to now() to avoid false assumption of future status.
// handled by ConvertIssueAssigneeHistory and ConvertIssueWithoutAssigneeHistory
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
