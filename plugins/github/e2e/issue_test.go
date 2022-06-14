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
			Owner: "panjf2000",
			Repo:  "ants",
			Config: models.Config{
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
	dataflowTester.ImportCsv("./raw_tables/_raw_github_api_issues.csv", "_raw_github_api_issues")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubIssue{})
	dataflowTester.Subtask(tasks.ExtractApiIssuesMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubIssue{}.TableName(),
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubIssue{}.TableName()),
		[]string{"github_id", "repo_id"},
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

	// import raw data table
	dataflowTester.ImportCsv("./raw_tables/_raw_github_api_pull_requests.csv", "_raw_github_api_pull_requests")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubPullRequest{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubPullRequest{}.TableName(),
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubPullRequest{}.TableName()),
		[]string{"github_id", "repo_id"},
		[]string{
			"number",
			"state",
			"title",
			"github_created_at",
			"github_updated_at",
			"closed_at",
			"additions",
			"deletions",
			"comments",
			"commits",
			"review_comments",
			"merged",
			"merged_at",
			"body",
			"type",
			"component",
			"merge_commit_sha",
			"head_ref",
			"base_ref",
			"base_commit_sha",
			"head_commit_sha",
			"url",
			"author_name",
			"author_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
