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
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestGithubJobsDataFlow(t *testing.T) {

	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
		},
		Repo: &models.GithubRepo{
			GithubId: 134018330,
		},
	}

	// import raw data table
	// SELECT * FROM _raw_github_api_jobs INTO OUTFILE "/tmp/_raw_github_api_jobs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_jobs.csv", "_raw_github_api_jobs")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubJob{})

	dataflowTester.Subtask(tasks.ExtractJobsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubJob{},
		"./snapshot_tables/_tool_github_jobs.csv",
		[]string{
			"connection_id",
			"repo_id",
			"id",
			"run_id",
			"run_url",
			"node_id",
			"head_sha",
			"url",
			"status",
			"conclusion",
			"started_at",
			"completed_at",
			"name",
			"steps",
			"check_run_url",
			"labels",
			"runner_id",
			"runner_name",
			"runner_group_id",
			"type",

			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}

func TestGithubRunsDataFlow(t *testing.T) {

	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
		},
		Repo: &models.GithubRepo{
			GithubId: 134018330,
		},
	}

	// import raw data table
	// SELECT * FROM _raw_github_api_runs INTO OUTFILE "/tmp/_raw_github_api_runs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_runs.csv", "_raw_github_api_runs")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubRun{})

	dataflowTester.Subtask(tasks.ExtractRunsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubRun{},
		"./snapshot_tables/_tool_github_runs.csv",
		[]string{
			"connection_id",
			"repo_id",
			"id",
			"name",
			"node_id",
			"head_branch",
			"head_sha",
			"path",
			"run_number",
			"event",
			"status",
			"conclusion",
			"workflow_id",
			"check_suite_id",
			"check_suite_node_id",
			"url",
			"html_url",
			"github_created_at",
			"github_updated_at",
			"run_attempt",
			"run_started_at",
			"jobs_url",
			"logs_url",
			"check_suite_url",
			"artifacts_url",
			"cancel_url",
			"rerun_url",
			"workflow_url",
			"type",

			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

}
