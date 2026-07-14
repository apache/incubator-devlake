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

type jiraIssue20260707 struct {
	FixVersions string `gorm:"type:text;column:fix_versions"`
}

func (jiraIssue20260707) TableName() string {
	return "_tool_jira_issues"
}

type changeFixVersionsToText20260707 struct{}

func (script *changeFixVersionsToText20260707) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &jiraIssue20260707{})
}

func (*changeFixVersionsToText20260707) Version() uint64 {
	return 20260707140000
}

func (*changeFixVersionsToText20260707) Name() string {
	return "change fix_versions type to text in _tool_jira_issues"
}
