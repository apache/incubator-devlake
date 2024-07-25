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
)

var _ plugin.MigrationScript = (*addTmpAccountIdToJiraIssueChangelogItem)(nil)

type jiraIssueChangelogItems20240709 struct {
	TmpFromAccountId string
	TmpToAccountId   string
}

func (jiraIssueChangelogItems20240709) TableName() string {
	return "_tool_jira_issue_changelog_items"
}

type addTmpAccountIdToJiraIssueChangelogItem struct{}

func (script *addTmpAccountIdToJiraIssueChangelogItem) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&jiraIssueChangelogItems20240709{})
}

func (*addTmpAccountIdToJiraIssueChangelogItem) Version() uint64 {
	return 20240709134200
}

func (*addTmpAccountIdToJiraIssueChangelogItem) Name() string {
	return "add TmpFromAccountId and TmpToAccountId to _tool_jira_issue_changelog_items"
}
