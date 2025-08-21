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

type jiraIssue20250619 struct {
	FixVersions string `gorm:"type:varchar(255)"`
}

func (jiraIssue20250619) TableName() string {
	return "_tool_jira_issues"
}

type addFixVersions20250619 struct{}

func (script *addFixVersions20250619) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &jiraIssue20250619{})
}

func (*addFixVersions20250619) Version() uint64 {
	return 20250619142316
}

func (*addFixVersions20250619) Name() string {
	return "add fix_versions field to _tool_jira_issues"
}
