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
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestPrDataFlow(t *testing.T) {
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
	dataflowTester.ImportCsv("./raw_tables/_raw_github_api_pull_requests.csv", "_raw_github_api_pull_requests")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubPullRequest{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubPullRequest{},
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

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubPullRequestLabel{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubPullRequestLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubPullRequestLabel{}.TableName()),
		[]string{"pull_id", "label_name"},
		[]string{},
	)

	// import raw data table
	dataflowTester.ImportCsv("./raw_tables/_raw_github_api_pull_request_commits.csv", "_raw_github_api_pull_request_commits")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubCommit{})
	dataflowTester.MigrateTableAndFlush(&models.GithubPullRequestCommit{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestCommitsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubCommit{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubCommit{}.TableName()),
		[]string{"sha"},
		[]string{
			"author_id",
			"author_name",
			"author_email",
			"authored_date",
			"committer_id",
			"committer_name",
			"committer_email",
			"committed_date",
			"message",
			"url",
			"additions",
			"deletions",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// import raw data table
	dataflowTester.ImportCsv("./raw_tables/_raw_github_api_pull_request_reviews.csv", "_raw_github_api_pull_request_reviews")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubReviewer{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestReviewsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		models.GithubReviewer{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubReviewer{}.TableName()),
		[]string{"github_id", "login", "pull_request_id"},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&code.PullRequestCommit{})
	dataflowTester.Subtask(tasks.ConvertPullRequestCommitsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		code.PullRequestCommit{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequestCommit{}.TableName()),
		[]string{"commit_sha", "pull_request_id"},
		[]string{},
	)

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertPullRequestsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		code.PullRequest{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequest{}.TableName()),
		[]string{"id"},
		[]string{
			"base_repo_id",
			"head_repo_id",
			"status",
			"title",
			"description",
			"url",
			"author_name",
			"author_id",
			"parent_pr_id",
			"pull_request_key",
			"created_date",
			"merged_date",
			"closed_date",
			"type",
			"component",
			"merge_commit_sha",
			"head_ref",
			"base_ref",
			"base_commit_sha",
			"head_commit_sha",
		},
	)

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&code.PullRequestLabel{})
	dataflowTester.Subtask(tasks.ConvertPullRequestLabelsMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		code.PullRequestLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequestLabel{}.TableName()),
		[]string{"pull_request_id", "label_name"},
		[]string{},
	)
}
