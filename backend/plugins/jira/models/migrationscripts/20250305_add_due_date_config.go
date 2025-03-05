package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type JiraScopeConfig20250305 struct {
	DueDateField string `mapstructure:"dueDateField,omitempty" json:"dueDateField" gorm:"type:varchar(255)"`
}

func (t JiraScopeConfig20250305) TableName() string {
	return "_tool_jira_scope_configs"
}

type updateScopeConfig struct{}

func (*updateScopeConfig) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&JiraScopeConfig20250305{},
	)
}

func (*updateScopeConfig) Version() uint64 {
	return 20250305092900
}

func (*updateScopeConfig) Name() string {
	return "jira update scope config, add due date field"
}