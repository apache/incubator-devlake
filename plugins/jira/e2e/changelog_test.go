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

func TestIssueChangelogDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// verify changelog conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_changelogs.csv", &models.JiraIssueChangelogs{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_changelog_items.csv", &models.JiraIssueChangelogItems{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_statuses_for_changelog.csv", &models.JiraStatus{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_changelog.csv", &models.JiraBoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertIssueChangelogsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueChangelogs{},
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
