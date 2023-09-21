package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addIncidentPriority struct {
	Priority string `gorm:"type:varchar(255)"`
}

func (*addIncidentPriority) TableName() string {
	return "_tool_pagerduty_incidents"
}
func (u *addIncidentPriority) Up(baseRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(baseRes,
		&addIncidentPriority{},
	)
}

func (*addIncidentPriority) Version() uint64 {
	return 20230920130004
}

func (*addIncidentPriority) Name() string {
	return "add priority to _tool_pagerduty_incidents table"
}
