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

func TestPrCommitDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_github_pull_requests.csv", models.GithubPullRequest{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_pull_request_commits.csv", "_raw_github_api_pull_request_commits")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubCommit{})
	dataflowTester.FlushTabler(&models.GithubPrCommit{})
	dataflowTester.FlushTabler(&models.GithubRepoCommit{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubCommit{},
		"./snapshot_tables/_tool_github_commits.csv",
		[]string{
			"sha",
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

	dataflowTester.VerifyTable(
		models.GithubPrCommit{},
		"./snapshot_tables/_tool_github_pull_request_commits.csv",
		[]string{
			"connection_id",
			"commit_sha",
			"pull_request_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify extraction
	dataflowTester.FlushTabler(&code.PullRequestCommit{})
	dataflowTester.Subtask(tasks.ConvertPullRequestCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestCommit{},
		"./snapshot_tables/pull_request_commits.csv",
		[]string{
			"commit_sha",
			"pull_request_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
