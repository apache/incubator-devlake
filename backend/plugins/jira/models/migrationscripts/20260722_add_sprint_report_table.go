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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type jiraSprintReport20260722 struct {
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

func (jiraSprintReport20260722) TableName() string {
	return "_tool_jira_sprint_reports"
}

type addSprintReportTable struct{}

func (script *addSprintReportTable) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &jiraSprintReport20260722{})
}

func (*addSprintReportTable) Version() uint64 {
	return 20260722000000
}

func (*addSprintReportTable) Name() string {
	return "add _tool_jira_sprint_reports table to persist Jira's frozen Sprint Report snapshot"
}
