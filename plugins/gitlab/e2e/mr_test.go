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
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabMrDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12955687,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_merge_requests.csv",
		"_raw_gitlab_api_merge_requests")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabMergeRequest{})
	dataflowTester.FlushTabler(&models.GitlabMrLabel{})
	dataflowTester.Subtask(tasks.ExtractApiMergeRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabMergeRequest{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GitlabMergeRequest{}.TableName()),
		[]string{"connection_id", "gitlab_id"},
		[]string{
			"iid",
			"project_id",
			"source_project_id",
			"target_project_id",
			"state",
			"title",
			"web_url",
			"user_notes_count",
			"work_in_progress",
			"source_branch",
			"target_branch",
			"merge_commit_sha",
			"merged_at",
			"gitlab_created_at",
			"closed_at",
			"merged_by_username",
			"description",
			"author_username",
			"author_user_id",
			"component",
			"first_comment_time",
			"review_rounds",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.GitlabMrLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GitlabMrLabel{}.TableName()),
		[]string{"connection_id", "mr_id", "label_name"},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.Subtask(tasks.ConvertApiMergeRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequest{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequest{}.TableName()),
		[]string{"id"},
		[]string{
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
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&code.PullRequestLabel{})
	dataflowTester.Subtask(tasks.ConvertMrLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequestLabel{}.TableName()),
		[]string{
			"pull_request_id",
			"label_name",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
