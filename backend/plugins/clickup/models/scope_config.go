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

// ClickUpScopeConfig allows a user to override how ClickUp raw statuses and
// task types map onto DevLake's standard domain values.
//
// Status: ClickUp statuses carry a `type` (open/unstarted/custom/done/closed)
// which the plugin maps automatically (open/unstarted -> TODO, custom ->
// IN_PROGRESS, done/closed -> DONE). When any of the IssueStatus* lists below
// are populated, a matching raw status name takes precedence over the
// type-derived default, so teams with bespoke workflows can classify custom
// statuses explicitly.
//
// Type: ClickUp has no native issue "type" on every task, so IssueType* are
// regular expressions matched against a task's derived type string. Precedence
// is INCIDENT > BUG > REQUIREMENT; a task matching none defaults to REQUIREMENT.
// A sensible default (bug -> BUG) is applied by the convertor when no config is
// set.
type ClickUpScopeConfig struct {
	common.ScopeConfig    `mapstructure:",squash" json:",inline" gorm:"embedded"`
	IssueStatusTodo       []string `mapstructure:"issueStatusTodo,omitempty" json:"issueStatusTodo" gorm:"type:json;serializer:json"`
	IssueStatusInProgress []string `mapstructure:"issueStatusInProgress,omitempty" json:"issueStatusInProgress" gorm:"type:json;serializer:json"`
	IssueStatusDone       []string `mapstructure:"issueStatusDone,omitempty" json:"issueStatusDone" gorm:"type:json;serializer:json"`
	IssueTypeRequirement  string   `mapstructure:"issueTypeRequirement,omitempty" json:"issueTypeRequirement" gorm:"type:varchar(255)"`
	IssueTypeBug          string   `mapstructure:"issueTypeBug,omitempty" json:"issueTypeBug" gorm:"type:varchar(255)"`
	IssueTypeIncident     string   `mapstructure:"issueTypeIncident,omitempty" json:"issueTypeIncident" gorm:"type:varchar(255)"`
}

func (ClickUpScopeConfig) TableName() string {
	return "_tool_clickup_scope_configs"
}

func (sc *ClickUpScopeConfig) SetConnectionId(c *ClickUpScopeConfig, connectionId uint64) {
	c.ConnectionId = connectionId
	c.ScopeConfig.ConnectionId = connectionId
}
