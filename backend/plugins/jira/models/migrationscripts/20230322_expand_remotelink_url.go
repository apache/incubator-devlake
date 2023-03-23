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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*expandRemotelinkUrl)(nil)

type jiraRemotelink20230322 struct {
	Url string
}

func (jiraRemotelink20230322) TableName() string {
	return "_tool_jira_remotelinks"
}

type jiraIssueCommit20230322 struct {
	CommitUrl string
}

func (jiraIssueCommit20230322) TableName() string {
	return "_tool_jira_issue_commits"
}

type expandRemotelinkUrl struct{}

func (script *expandRemotelinkUrl) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	// expand _tool_jira_remotelinks.url to LONGTEXT
	err := migrationhelper.ChangeColumnsType[jiraRemotelink20230322](
		basicRes,
		script,
		jiraRemotelink20230322{}.TableName(),
		[]string{"url"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&jiraRemotelink20230322{},
				"url",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? is not null ", tmpColumnParams...),
			)
		},
	)
	if err != nil {
		return err
	}
	// expand _tool_jira_issue_commits.commit_url to LONGTEXT
	err = migrationhelper.ChangeColumnsType[jiraIssueCommit20230322](
		basicRes,
		script,
		jiraIssueCommit20230322{}.TableName(),
		[]string{"commit_url"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&jiraIssueCommit20230322{},
				"commit_url",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? is not null ", tmpColumnParams...),
			)
		},
	)
	return err
}

func (*expandRemotelinkUrl) Version() uint64 {
	return 20230322153324
}

func (*expandRemotelinkUrl) Name() string {
	return "expand _tool_jira_remotelinks.url and _tool_jira_issue_commits.commit_url to LONGTEXT"
}
