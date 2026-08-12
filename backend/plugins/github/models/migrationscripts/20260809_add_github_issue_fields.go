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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type githubIssueFieldValue20260809 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueId      int    `gorm:"primaryKey"`
	FieldId      int    `gorm:"primaryKey"`
	FieldName    string `gorm:"type:varchar(255);index"`
	DataType     string `gorm:"type:varchar(50)"`
	Value        string `gorm:"type:text"`
	RawValue     string `gorm:"type:text"`
	OptionColor  string `gorm:"type:varchar(50)"`
	archived.NoPKModel
}

func (githubIssueFieldValue20260809) TableName() string {
	return "_tool_github_issue_field_values"
}

type githubScopeConfig20260809 struct {
	IssueFieldPriority   string `gorm:"type:varchar(255)"`
	IssueFieldSeverity   string `gorm:"type:varchar(255)"`
	IssueFieldComponent  string `gorm:"type:varchar(255)"`
	IssueFieldStoryPoint string `gorm:"type:varchar(255)"`
	IssueFieldDueDate    string `gorm:"type:varchar(255)"`
}

func (githubScopeConfig20260809) TableName() string {
	return "_tool_github_scope_configs"
}

type addGithubIssueFields struct{}

func (*addGithubIssueFields) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&githubIssueFieldValue20260809{},
		&githubScopeConfig20260809{},
	)
}

func (*addGithubIssueFields) Version() uint64 {
	return 20260809000001
}

func (*addGithubIssueFields) Name() string {
	return "add github issue field values table and scope config issue field mappings"
}
