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
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/dora/impl"
	"github.com/apache/incubator-devlake/plugins/dora/tasks"
)

func TestCalculateCLTimeDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			RepoId: "github:GithubRepo:1:384111310",
			TransformationRules: tasks.TransformationRules{
				ProductionPattern: "(?i)deploy",
			},
		},
	}

	dataflowTester.FlushTabler(&code.PullRequest{})

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/lake_cicd_tasks.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_requests.csv", &code.PullRequest{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_request_comments.csv", &code.PullRequestComment{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_request_commits.csv", &code.PullRequestCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/commits.csv", &code.Commit{})

	// verify converter
	dataflowTester.Subtask(tasks.CalculateChangeLeadTimeMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequest{},
		"./snapshot_tables/pull_requests.csv",
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
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
			"coding_timespan",
			"review_lag",
			"review_timespan",
			"deploy_timespan",
			"change_timespan",
			"orig_coding_timespan",
			"orig_review_lag",
			"orig_review_timespan",
			"orig_deploy_timespan",
		},
	)
}
