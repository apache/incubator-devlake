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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
)

func TestWorklogDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_worklogs.csv", "_raw_jira_api_worklogs")

	// verify worklog extraction
	dataflowTester.FlushTabler(&models.JiraWorklog{})
	dataflowTester.Subtask(tasks.ExtractWorklogsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraWorklog{},
		"./snapshot_tables/_tool_jira_worklogs.csv",
		[]string{"connection_id", "issue_id", "worklog_id"},
		[]string{
			"author_id",
			"update_author_id",
			"time_spent",
			"time_spent_seconds",
			"updated",
			"started",
			"issue_updated",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify worklog conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_worklog.csv", &models.JiraBoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueWorklog{})
	dataflowTester.Subtask(tasks.ConvertWorklogsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueWorklog{},
		"./snapshot_tables/worklogs.csv",
		[]string{"id"},
		[]string{
			"author_id",
			"comment",
			"time_spent_minutes",
			"logged_date",
			"started_date",
			"issue_id",
		},
	)
}
