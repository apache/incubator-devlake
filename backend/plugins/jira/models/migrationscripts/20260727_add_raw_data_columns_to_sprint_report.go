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
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/migrationscripts/archived"
	"github.com/apache/devlake/helpers/migrationhelper"
)

// jiraSprintReport20260727 mirrors the JiraSprintReport model. The original
// migration (20260722) created _tool_jira_sprint_reports without embedding
// common.NoPKModel, so the _raw_data_table / _raw_data_params / _raw_data_id /
// _raw_data_remark columns (and created_at / updated_at) were missing. The
// runtime model expects them, which made the ApiExtractor's cleanup query
// (WHERE _raw_data_table = ? AND _raw_data_params = ?) fail with
// "Unknown column '_raw_data_table' in 'where clause'". Re-running
// AutoMigrateTables adds the missing columns without dropping existing data.
type jiraSprintReport20260727 struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`
	SprintId     uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`

	IssueKey                 string `gorm:"type:varchar(255)"`
	Bucket                   string `gorm:"type:varchar(32);index"`
	Done                     bool
	StoryPointsAtSprintStart *float64
	StoryPointsAtSprintEnd   *float64
}

func (jiraSprintReport20260727) TableName() string {
	return "_tool_jira_sprint_reports"
}

type addRawDataColumnsToSprintReport struct{}

func (script *addRawDataColumnsToSprintReport) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &jiraSprintReport20260727{})
}

func (*addRawDataColumnsToSprintReport) Version() uint64 {
	return 20260727000000
}

func (*addRawDataColumnsToSprintReport) Name() string {
	return "add missing _raw_data_* columns to _tool_jira_sprint_reports"
}
