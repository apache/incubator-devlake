package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// IssueLeadTimeMetric stores the calculated lead time from the first 'In Progress' status
// to the first 'Done' status for a specific issue.
type IssueLeadTimeMetric struct {
	common.NoPKModel `json:"-" gorm:"primaryKey;autoIncrement:false"`

	ProjectName string `json:"projectName" gorm:"primaryKey;type:varchar(255)"`
	IssueId     string `json:"issueId" gorm:"primaryKey;type:varchar(255)"`

	FirstInProgressDate *time.Time `json:"firstInProgressDate"`
	FirstDoneDate       *time.Time `json:"firstDoneDate"`

	// Lead time in minutes from first 'In Progress' to first 'Done'
	InProgressToDoneMinutes *int64 `json:"inProgressToDoneMinutes"`
}

// TableName specifies the database table name for this model.
func (IssueLeadTimeMetric) TableName() string {
	return "_tool_dora_issue_lead_time_metrics"
}
