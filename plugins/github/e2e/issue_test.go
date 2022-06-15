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
	"fmt"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestIssueDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", plugin)

	githubRepository := &models.GithubRepo{
		GithubId: 134018330,
	}
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			Config: models.Config{
				PrType:               "type/(.*)$",
				PrComponent:          "component/(.*)$",
				PrBodyClosePattern:   "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)",
				IssueSeverity:        "severity/(.*)$",
				IssuePriority:        "^(highest|high|medium|low)$",
				IssueComponent:       "component/(.*)$",
				IssueTypeBug:         "^(bug|failure|error)$",
				IssueTypeIncident:    "",
				IssueTypeRequirement: "^(feat|feature|proposal|requirement)$",
			},
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_issues.csv", "_raw_github_api_issues")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubIssue{})
	dataflowTester.FlushTabler(&models.GithubIssueLabel{})
	dataflowTester.Subtask(tasks.ExtractApiIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubIssue{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubIssue{}.TableName()),
		[]string{"connection_id", "github_id", "repo_id"},
		[]string{
			"number",
			"state",
			"title",
			"body",
			"priority",
			"type",
			"status",
			"author_id",
			"author_name",
			"assignee_id",
			"assignee_name",
			"lead_time_minutes",
			"url",
			"closed_at",
			"github_created_at",
			"github_updated_at",
			"severity",
			"component",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	//// verify extraction
	//dataflowTester.VerifyTable(
	//	models.GithubIssueLabel{},
	//	fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubIssueLabel{}.TableName()),
	//	[]string{"issue_id", "label_name"},
	//	[]string{
	//		"_raw_data_params",
	//		"_raw_data_table",
	//		"_raw_data_id",
	//		"_raw_data_remark",
	//	},
	//)
	//
	//// verify extraction
	//dataflowTester.FlushTabler(&ticket.Issue{})
	//dataflowTester.FlushTabler(&ticket.BoardIssue{})
	//dataflowTester.Subtask(tasks.ConvertIssuesMeta, taskData)
	//dataflowTester.VerifyTable(
	//	ticket.Issue{},
	//	fmt.Sprintf("./snapshot_tables/%s.csv", ticket.Issue{}.TableName()),
	//	[]string{"id"},
	//	[]string{
	//		"url",
	//		"icon_url",
	//		"issue_key",
	//		"title",
	//		"description",
	//		"epic_key",
	//		"type",
	//		"status",
	//		"original_status",
	//		"story_point",
	//		"resolution_date",
	//		"created_date",
	//		"updated_date",
	//		"lead_time_minutes",
	//		"parent_issue_id",
	//		"priority",
	//		"original_estimate_minutes",
	//		"time_spent_minutes",
	//		"time_remaining_minutes",
	//		"creator_id",
	//		"creator_name",
	//		"assignee_id",
	//		"assignee_name",
	//		"severity",
	//		"component",
	//	},
	//)
	//dataflowTester.VerifyTable(
	//	ticket.BoardIssue{},
	//	fmt.Sprintf("./snapshot_tables/%s.csv", ticket.BoardIssue{}.TableName()),
	//	[]string{"board_id", "issue_id"},
	//	[]string{},
	//)
	//
	//// verify extraction
	//dataflowTester.FlushTabler(&ticket.IssueLabel{})
	//dataflowTester.Subtask(tasks.ConvertIssueLabelsMeta, taskData)
	//dataflowTester.VerifyTable(
	//	ticket.IssueLabel{},
	//	fmt.Sprintf("./snapshot_tables/%s.csv", ticket.IssueLabel{}.TableName()),
	//	[]string{"issue_id", "label_name"},
	//	[]string{},
	//)
}
