package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type TapdScopeConfig20250305 struct {
	BugDueDateField   string `mapstructure:"bugDueDateField,omitempty" json:"bugDueDateField" gorm:"column:bug_due_date_field"`
	TaskDueDateField  string `mapstructure:"taskDueDateField,omitempty" json:"taskDueDateField" gorm:"column:task_due_date_field"`
	StoryDueDateField string `mapstructure:"storyDueDateField,omitempty" json:"storyDueDateField" gorm:"column:story_due_date_field"`
}

func (t TapdScopeConfig20250305) TableName() string {
	return "_tool_tapd_scope_configs"
}

type updateScopeConfig20250305 struct{}

func (u updateScopeConfig20250305) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&TapdScopeConfig20250305{},
	)
}

func (u updateScopeConfig20250305) Version() uint64 {
	return 20250305093400
}

func (u updateScopeConfig20250305) Name() string {
	return "tapd update scope config, add due date fields"
}