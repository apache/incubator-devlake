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
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestPrDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)

	githubRepository := &models.GithubRepo{
		GithubId: 134018330,
	}
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			TransformationRules: &models.TransformationRules{
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
	dataflowTester.FlushTabler(&models.GithubPrLabel{})
	dataflowTester.FlushTabler(&models.GithubAccount{})
	dataflowTester.FlushTabler(&models.GithubRepoAccount{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubPullRequest{},
		"./snapshot_tables/_tool_github_pull_requests.csv",
		[]string{
			"connection_id",
			"github_id",
			"repo_id",
			"head_repo_id",
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
		models.GithubPrLabel{},
		"./snapshot_tables/_tool_github_pull_request_labels.csv",
		[]string{"connection_id", "pull_id", "label_name"},
	)

	dataflowTester.VerifyTable(
		models.GithubRepoAccount{},
		"./snapshot_tables/_tool_github_accounts_in_pr.csv",
		[]string{
			"connection_id",
			"account_id",
			"repo_github_id",
			"login",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify pr conversion
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequest{},
		"./snapshot_tables/pull_requests.csv",
		[]string{
			"id",
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
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify label conversion
	dataflowTester.FlushTabler(&code.PullRequestLabel{})
	dataflowTester.Subtask(tasks.ConvertPullRequestLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestLabel{},
		"./snapshot_tables/pull_request_labels.csv",
		[]string{
			"pull_request_id",
			"label_name",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
