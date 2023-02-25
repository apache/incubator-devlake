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

package e2e

import (
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"testing"
)

func TestConvertIssueRepoCommitsDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
			TransformationRules: &tasks.JiraTransformationRule{
				RemotelinkCommitShaPattern: `.*/commit/(.*)`,
				RemotelinkRepoPattern: []string{
					`https://bitbucket.org/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commits/(?P<commit_sha>\w{40})`,
					`https://gitlab.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commit/(?P<commit_sha>\w{40})`,
					`https://github.com/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commit/(?P<commit_sha>\w{40})`,
				},
			},
		},
	}
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_commits_for_ConvertIssueRepoCommits.csv", &models.JiraIssueCommit{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_ConvertIssueRepoCommits.csv", &models.JiraBoardIssue{})
	dataflowTester.Subtask(tasks.ConvertIssueRepoCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.IssueRepoCommit{},
		"./snapshot_tables/issue_repo_commits.csv",
		e2ehelper.ColumnWithRawData(
			"issue_id",
			"repo_url",
			"commit_sha",
			"host",
			"namespace",
			"repo_name",
		),
	)
}
