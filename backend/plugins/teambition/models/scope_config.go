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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type TeambitionScopeConfig struct {
	common.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	TypeMappings       map[string]string `mapstructure:"typeMappings,omitempty" json:"typeMappings" gorm:"serializer:json"`
	StatusMappings     map[string]string `mapstructure:"statusMappings,omitempty" json:"statusMappings" gorm:"serializer:json"`
	BugDueDateField    string            `mapstructure:"bugDueDateField,omitempty" json:"bugDueDateField" gorm:"column:bug_due_date_field"`
	TaskDueDateField   string            `mapstructure:"taskDueDateField,omitempty" json:"taskDueDateField" gorm:"column:task_due_date_field"`
	StoryDueDateField  string            `mapstructure:"storyDueDateField,omitempty" json:"storyDueDateField" gorm:"column:story_due_date_field"`
}

func (t TeambitionScopeConfig) TableName() string {
	return "_tool_teambition_scope_configs"
}

func (t *TeambitionScopeConfig) SetConnectionId(c *TeambitionScopeConfig, connectionId uint64) {
	c.ConnectionId = connectionId
	c.ScopeConfig.ConnectionId = connectionId
}
