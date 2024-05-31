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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabMrCommitDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12345678,
			ScopeConfig:  new(models.GitlabScopeConfig),
		},
	}

	// Prepare _tool_gitlab_pull_requests for mr_commit convertor test
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_merge_requests_for_mr_commit_test.csv",
		"_raw_gitlab_api_merge_requests")
	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabMergeRequest{})
	dataflowTester.FlushTabler(&models.GitlabMrLabel{})
	dataflowTester.FlushTabler(&models.GitlabAssignee{})
	dataflowTester.FlushTabler(&models.GitlabReviewer{})
	dataflowTester.Subtask(tasks.ExtractApiMergeRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabMergeRequest{},
		"./snapshot_tables/_tool_gitlab_merge_requests_for_mr_commit_test.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"gitlab_id",
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
		),
	)

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_merge_request_commits.csv",
		"_raw_gitlab_api_merge_request_commits")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabCommit{})
	dataflowTester.FlushTabler(&models.GitlabProjectCommit{})
	dataflowTester.FlushTabler(&models.GitlabMrCommit{})
	dataflowTester.Subtask(tasks.ExtractApiMrCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabCommit{},
		"./snapshot_tables/_tool_gitlab_commits.csv",
		e2ehelper.ColumnWithRawData(
			"sha",
			"title",
			"message",
			"short_id",
			"author_name",
			"author_email",
			"authored_date",
			"committer_name",
			"committer_email",
			"committed_date",
			"web_url",
			"additions",
			"deletions",
			"total",
		),
	)
	dataflowTester.VerifyTable(
		models.GitlabProjectCommit{},
		"./snapshot_tables/_tool_gitlab_project_commits.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"gitlab_project_id",
			"commit_sha",
		),
	)

	dataflowTester.VerifyTableWithOptions(models.GitlabMrCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_gitlab_mr_commits.csv",
		IgnoreTypes: []interface{}{common.Model{}},
	})

	// verify conversion
	dataflowTester.FlushTabler(&code.PullRequestCommit{})
	dataflowTester.Subtask(tasks.ConvertApiMrCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(code.PullRequestCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/pull_request_commits.csv",
		IgnoreTypes: []interface{}{common.Model{}},
	})

	// verify conversion
	dataflowTester.FlushTabler(&code.Commit{})
	dataflowTester.FlushTabler(&code.RepoCommit{})
	dataflowTester.Subtask(tasks.ConvertCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		code.Commit{},
		"./snapshot_tables/commits.csv",
		e2ehelper.ColumnWithRawData(
			"sha",
			"additions",
			"deletions",
			"dev_eq",
			"message",
			"author_name",
			"author_email",
			"authored_date",
			"author_id",
			"committer_name",
			"committer_email",
			"committed_date",
			"committer_id",
		),
	)

	dataflowTester.VerifyTable(
		code.RepoCommit{},
		"./snapshot_tables/repo_commits.csv",
		e2ehelper.ColumnWithRawData(
			"repo_id",
			"commit_sha",
		),
	)
}
