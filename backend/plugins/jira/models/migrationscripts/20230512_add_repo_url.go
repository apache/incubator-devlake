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

type JiraIssueCommit20230512 struct {
	RepoUrl string
}

func (JiraIssueCommit20230512) TableName() string {
	return "_tool_jira_issue_commits"
}

type addRepoUrl struct{}

func (script *addRepoUrl) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &JiraIssueCommit20230512{})
}

func (*addRepoUrl) Version() uint64 {
	return 20230512113738
}

func (*addRepoUrl) Name() string {
	return "add repo_url to _tool_jira_issue_commits"
}
