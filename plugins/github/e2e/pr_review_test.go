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
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestPrReviewDataFlow(t *testing.T) {
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
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_github_pull_requests.csv", models.GithubPullRequest{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_pull_request_reviews.csv", "_raw_github_api_pull_request_reviews")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubReviewer{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestReviewersMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubReviewer{},
		"./snapshot_tables/_tool_github_reviewers.csv",
		[]string{"connection_id", "github_id", "pull_request_id"},
		[]string{
			"login",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
