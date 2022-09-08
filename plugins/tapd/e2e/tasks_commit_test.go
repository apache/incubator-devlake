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

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdTaskCommitDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}

	// task status
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_task_commits.csv",
		"_raw_tapd_api_task_commits")
	// verify extraction
	dataflowTester.FlushTabler(&models.TapdTaskCommit{})
	dataflowTester.Subtask(tasks.ExtractTaskCommitMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdTaskCommit{},
		"./snapshot_tables/_tool_tapd_task_commits.csv",
		[]string{
			"connection_id",
			"id",
			"user_id",
			"hook_user_name",
			"commit_id",
			"workspace_id",
			"message",
			"path",
			"web_url",
			"hook_project_name",
			"ref",
			"ref_status",
			"git_env",
			"file_commit",
			"commit_time",
			"created",
			"task_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.FlushTabler(&crossdomain.IssueCommit{})
	dataflowTester.Subtask(tasks.ConvertTaskCommitMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.IssueCommit{},
		"./snapshot_tables/issue_commits_task.csv",
		[]string{
			"issue_id",
			"commit_sha",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
