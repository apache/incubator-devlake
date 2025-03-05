package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type ZentaoScopeConfig20250305 struct {
	BugDueDateField   string `mapstructure:"bugDueDateField,omitempty" json:"bugDueDateField"`
	TaskDueDateField  string `mapstructure:"taskDueDateField,omitempty" json:"taskDueDateField"`
	StoryDueDateField string `mapstructure:"storyDueDateField,omitempty" json:"storyDueDateField"`
}

func (t ZentaoScopeConfig20250305) TableName() string {
	return "_tool_zentao_scope_configs"
}

type updateScopeConfig struct{}

func (*updateScopeConfig) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&ZentaoScopeConfig20250305{},
	)
}

func (*updateScopeConfig) Version() uint64 {
	return 20250305092300
}

func (*updateScopeConfig) Name() string {
	return "zentao update scope config, add due date fields"
}
