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

func TestChangelogDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issue_changelogs.csv", "_raw_jira_api_issue_changelogs")

	// verify changelog extraction
	dataflowTester.FlushTabler(&models.JiraIssueChangelogs{})
	dataflowTester.FlushTabler(&models.JiraIssueChangelogItems{})
	dataflowTester.FlushTabler(&models.JiraAccount{})
	dataflowTester.Subtask(tasks.ExtractIssueChangelogsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraIssueChangelogs{},
		"./snapshot_tables/_tool_jira_issue_changelogs.csv",
		[]string{
			"connection_id",
			"changelog_id",
			"issue_id",
			"author_account_id",
			"author_display_name",
			"author_active",
			"created",
			"issue_updated",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.VerifyTable(
		models.JiraIssueChangelogItems{},
		"./snapshot_tables/_tool_jira_issue_changelog_items.csv",
		[]string{
			"connection_id",
			"changelog_id",
			"field",
			"field_type",
			"field_id",
			"from_value",
			"from_string",
			"to_value",
			"to_string",
		},
	)

	// verify changelog conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_changelog.csv", &models.JiraBoardIssue{})
	dataflowTester.FlushTabler(&ticket.Changelog{})
	dataflowTester.Subtask(tasks.ConvertIssueChangelogsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Changelog{},
		"./snapshot_tables/issue_changelogs.csv",
		[]string{
			"id",
			"issue_id",
			"author_id",
			"author_name",
			"field_id",
			"field_name",
			"original_from_value",
			"original_to_value",
			"from_value",
			"to_value",
			"created_date",
		},
	)
}
