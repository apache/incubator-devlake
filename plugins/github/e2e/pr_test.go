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
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			TransformationRules: models.TransformationRules{
				PrType:             "type/(.*)$",
				PrComponent:        "component/(.*)$",
				PrBodyClosePattern: "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)",
			},
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_pull_requests.csv", "_raw_github_api_pull_requests")

	// verify pr extraction
	dataflowTester.FlushTabler(&models.GithubPullRequest{})
	dataflowTester.FlushTabler(&models.GithubPullRequestLabel{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubPullRequest{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubPullRequest{}.TableName()),
		[]string{"connection_id", "github_id", "repo_id"},
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

	dataflowTester.VerifyTable(
		models.GithubPullRequestLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubPullRequestLabel{}.TableName()),
		[]string{"connection_id", "pull_id", "label_name"},
		[]string{},
	)

	// verify pr conversion
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
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

	// verify label conversion
	dataflowTester.FlushTabler(&code.PullRequestLabel{})
	dataflowTester.Subtask(tasks.ConvertPullRequestLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequestLabel{}.TableName()),
		[]string{"pull_request_id", "label_name"},
		[]string{},
	)
}
