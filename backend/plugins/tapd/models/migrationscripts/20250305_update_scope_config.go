/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrationscripts

import (
	"time"

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

type TapdBug20250310 struct {
	DueDate *time.Time
}

func (t TapdBug20250310) TableName() string {
	return "_tool_tapd_bugs"
}

type TapdTask20250310 struct {
	DueDate *time.Time
}

func (t TapdTask20250310) TableName() string {
	return "_tool_tapd_tasks"
}

type TapdStory20250310 struct {
	DueDate *time.Time
}

func (t TapdStory20250310) TableName() string {
	return "_tool_tapd_stories"
}

type updateScopeConfig20250305 struct{}

func (u updateScopeConfig20250305) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&TapdScopeConfig20250305{},
		&TapdBug20250310{},
		&TapdStory20250310{},
		&TapdTask20250310{},
	)
}

func (u updateScopeConfig20250305) Version() uint64 {
	return 20250305093400
}

func (u updateScopeConfig20250305) Name() string {
	return "tapd update scope config, add due date fields"
}
