package models

import (
	"time"
)

// IssueLeadTimeMetric tracks lead time for issues from in-progress to done
type IssueLeadTimeMetric struct {
	ProjectName    string     `json:"projectName" gorm:"primaryKey;type:varchar(255)"`
	IssueId        string     `json:"issueId" gorm:"primaryKey;type:varchar(255)"`
	InProgressDate *time.Time `json:"InProgressDate"`
	DoneDate       *time.Time `json:"DoneDate"`

	// Lead time in minutes from first 'In Progress' to first 'Done'
	InProgressToDoneMinutes *int64 `json:"inProgressToDoneMinutes"`
}

// TableName specifies the database table name
func (IssueLeadTimeMetric) TableName() string {
	return "_tool_dora_issue_lead_time_metrics"
}
