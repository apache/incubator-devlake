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

	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestPrReviewDataFlow(t *testing.T) {
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
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_pull_request_reviews.csv", "_raw_github_api_pull_request_reviews")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubReviewer{})
	dataflowTester.FlushTabler(&models.GithubPrReview{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestReviewsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubPrReview{},
		"./snapshot_tables/_tool_github_pull_request_reviews.csv",
		[]string{
			"connection_id",
			"github_id",
			"pull_request_id",
			"body",
			"author_username",
			"author_user_id",
			"github_submit_at",
			"commit_sha",
			"state",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.GithubReviewer{},
		"./snapshot_tables/_tool_github_reviewers.csv",
		[]string{
			"connection_id",
			"reviewer_id",
			"pull_request_id",
			"username",
			"name",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.FlushTabler(&code.PullRequestComment{})
	dataflowTester.Subtask(tasks.ConvertPullRequestReviewsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestComment{},
		"./snapshot_tables/pull_request_review_comments.csv",
		[]string{
			"pull_request_id",
			"body",
			"account_id",
			"created_date",
			"commit_sha",
			"type",
			"review_id",
			"status",
		},
	)
}
