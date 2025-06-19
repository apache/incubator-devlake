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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addAppIdBack)(nil)

type teambitionConnection20250529 struct {
	AppId     string
	SecretKey string `gorm:"serializer:encdec"`
}

func (teambitionConnection20250529) TableName() string {
	return "_tool_teambition_connections"
}

type teambitionScopeConfig20250529 struct {
	Entities          []string          `gorm:"type:json;serializer:json" json:"entities" mapstructure:"entities"`
	ConnectionId      uint64            `json:"connectionId" gorm:"index" validate:"required" mapstructure:"connectionId,omitempty"`
	Name              string            `mapstructure:"name" json:"name" gorm:"type:varchar(255);uniqueIndex" validate:"required"`
	TypeMappings      map[string]string `mapstructure:"typeMappings,omitempty" json:"typeMappings" gorm:"serializer:json"`
	StatusMappings    map[string]string `mapstructure:"statusMappings,omitempty" json:"statusMappings" gorm:"serializer:json"`
	BugDueDateField   string            `mapstructure:"bugDueDateField,omitempty" json:"bugDueDateField" gorm:"column:bug_due_date_field"`
	TaskDueDateField  string            `mapstructure:"taskDueDateField,omitempty" json:"taskDueDateField" gorm:"column:task_due_date_field"`
	StoryDueDateField string            `mapstructure:"storyDueDateField,omitempty" json:"storyDueDateField" gorm:"column:story_due_date_field"`
}

func (t teambitionScopeConfig20250529) TableName() string {
	return "_tool_teambition_scope_configs"
}

type addAppIdBack struct{}

func (*addAppIdBack) Up(basicRes context.BasicRes) errors.Error {
	basicRes.GetLogger().Warn(nil, "*********")
	err := migrationhelper.AutoMigrateTables(basicRes, &teambitionConnection20250529{})
	basicRes.GetLogger().Warn(err, "err connection")
	err = migrationhelper.AutoMigrateTables(basicRes, &teambitionScopeConfig20250529{})
	basicRes.GetLogger().Warn(err, "err scope")
	return err
}

func (*addAppIdBack) Version() uint64 {
	return 20250529165745
}

func (*addAppIdBack) Name() string {
	return "add app id back to teambition_connections"
}
