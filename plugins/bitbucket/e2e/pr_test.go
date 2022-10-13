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

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestPrDataFlow(t *testing.T) {
	var plugin impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", plugin)

	bitbucketRepository := &models.BitbucketRepo{
		BitbucketId: "panjf2000/ants",
	}
	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			TransformationRules: models.TransformationRules{
				PrType:             "type/(.*)$",
				PrComponent:        "component/(.*)$",
				PrBodyClosePattern: "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/api.bitbucket.org\\/2.0\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)",
			},
		},
		Repo: bitbucketRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_api_pull_requests.csv", "_raw_bitbucket_api_pull_requests")

	// verify pr extraction
	dataflowTester.FlushTabler(&models.BitbucketPullRequest{})
	dataflowTester.FlushTabler(&models.BitbucketAccount{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketPullRequest{},
		"./snapshot_tables/_tool_bitbucket_pull_requests.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"bitbucket_id",
			"repo_id",
			"number",
			"state",
			"title",
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
		models.BitbucketAccount{},
		"./snapshot_tables/_tool_bitbucket_accounts_in_pr.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"user_name",
			"account_id",
			"account_status",
			"display_name",
			"avatar_url",
			"html_url",
			"uuid",
			"has2_fa_enabled",
		),
	)

	// verify pr conversion
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertPullRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequest{},
		"./snapshot_tables/pull_requests.csv",
		e2ehelper.ColumnWithRawData(
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
			"merged_date",
			"closed_date",
			"type",
			"component",
			"merge_commit_sha",
			"head_ref",
			"base_ref",
			"base_commit_sha",
			"head_commit_sha",
		),
	)

}
