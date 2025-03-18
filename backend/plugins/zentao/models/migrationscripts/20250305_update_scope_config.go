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

type ZentaoScopeConfig20250305 struct {
	BugDueDateField   string `mapstructure:"bugDueDateField,omitempty" json:"bugDueDateField"`
	TaskDueDateField  string `mapstructure:"taskDueDateField,omitempty" json:"taskDueDateField"`
	StoryDueDateField string `mapstructure:"storyDueDateField,omitempty" json:"storyDueDateField"`
}

func (t ZentaoScopeConfig20250305) TableName() string {
	return "_tool_zentao_scope_configs"
}

type ZentaoBug20250310 struct {
	DueDate *time.Time
}

func (t ZentaoBug20250310) TableName() string {
	return "_tool_zentao_bugs"
}

type ZentaoTask20250310 struct {
	DueDate *time.Time
}

func (t ZentaoTask20250310) TableName() string {
	return "_tool_zentao_tasks"
}

type ZentaoStory20250310 struct {
	DueDate *time.Time
}

func (t ZentaoStory20250310) TableName() string {
	return "_tool_zentao_stories"
}

type updateScopeConfig struct{}

func (*updateScopeConfig) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&ZentaoScopeConfig20250305{},
		&ZentaoBug20250310{},
		&ZentaoTask20250310{},
		&ZentaoStory20250310{},
	)
}

func (*updateScopeConfig) Version() uint64 {
	return 20250305092300
}

func (*updateScopeConfig) Name() string {
	return "zentao update scope config, add due date fields"
}
