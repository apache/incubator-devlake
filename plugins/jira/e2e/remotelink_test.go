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
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
)

func TestRemotelinkDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId:        2,
			BoardId:             8,
			TransformationRules: tasks.TransformationRules{RemotelinkCommitShaPattern: ".*/commit/(.*)"},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_remotelinks.csv", "_raw_jira_api_remotelinks")

	// verify remotelink extraction
	dataflowTester.FlushTabler(&models.JiraRemotelink{})
	dataflowTester.FlushTabler(&models.JiraIssueCommit{})
	dataflowTester.Subtask(tasks.ExtractRemotelinksMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraRemotelink{},
		"./snapshot_tables/_tool_jira_remotelinks.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"remotelink_id",
			"issue_id",
			"self",
			"title",
			"url",
			"issue_updated",
		),
	)
	dataflowTester.VerifyTable(
		models.JiraIssueCommit{},
		"./snapshot_tables/_tool_jira_issue_commits.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"issue_id",
			"commit_sha",
			"commit_url",
		),
	)
}
