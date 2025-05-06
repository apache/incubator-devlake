package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

// Define the actual table structure directly in the migration script
type issueLeadTimeMetric struct {
	ProjectName             string `gorm:"primaryKey;type:varchar(255)"`
	IssueId                 string `gorm:"primaryKey;type:varchar(255)"`
	InProgressDate          *time.Time
	DoneDate                *time.Time
	InProgressToDoneMinutes *int64
}

// TableName specifies the table name
func (issueLeadTimeMetric) TableName() string {
	return "_tool_dora_issue_lead_time_metrics"
}

type addIssueLeadTimeMetricsTable struct{}

func (script *addIssueLeadTimeMetricsTable) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	// Use our directly defined model instead of importing from models
	return db.AutoMigrate(&issueLeadTimeMetric{})
}

func (*addIssueLeadTimeMetricsTable) Version() uint64 {
	return 2025042401
}

func (*addIssueLeadTimeMetricsTable) Name() string {
	return "dora add _tool_dora_issue_lead_time_metrics table"
}
