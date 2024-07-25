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
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/tasks"
)

func TestPrDataFlow(t *testing.T) {
	var plugin impl.BitbucketServer
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket_server", plugin)

	taskData := &tasks.BitbucketServerTaskData{
		Options: &tasks.BitbucketServerOptions{
			ConnectionId: 3,
			FullName:     "TP/repos/first-repo",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_server_api_pull_requests.csv", "_raw_bitbucket_server_api_pull_requests")

	// verify pr extraction
	dataflowTester.FlushTabler(&models.BitbucketServerPullRequest{})
	dataflowTester.FlushTabler(&models.BitbucketServerUser{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketServerPullRequest{},
		"./snapshot_tables/_tool_bitbucket_server_pull_requests.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"repo_id",
			"bitbucket_id",
			"number",
			"base_repo_id",
			"head_repo_id",
			"state",
			"title",
			"description",
			"closed_at",
			"comment_count",
			"commits",
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
		),
	)

	dataflowTester.VerifyTable(
		models.BitbucketServerUser{},
		"./snapshot_tables/_tool_bitbucket_server_users_in_pr.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"bitbucket_id",
			"name",
			"email_address",
			"active",
			"slug",
			"type",
			"account_status",
			"display_name",
			"html_url",
		),
	)

	// verify pr conversion
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequest{},
		"./snapshot_tables/pull_requests.csv",
		e2ehelper.ColumnWithRawData(
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
			"original_status",
		),
	)
}
