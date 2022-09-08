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
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabJobDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    44,
		},
	}
	// import raw data table
	// SELECT * FROM _raw_gitlab_api_job INTO OUTFILE "/tmp/_raw_gitlab_api_job.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_job.csv", "_raw_gitlab_api_job")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabJob{})
	dataflowTester.Subtask(tasks.ExtractApiJobsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabJob{},
		"./snapshot_tables/_tool_gitlab_jobs.csv",
		[]string{
			"connection_id",
			"gitlab_id",
			"project_id",
			"pipeline_id",
			"status",
			"stage",
			"name",
			"ref",
			"tag",
			"allow_failure",
			"duration",
			"web_url",

			"gitlab_created_at",
			"started_at",
			"finished_at",

			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
