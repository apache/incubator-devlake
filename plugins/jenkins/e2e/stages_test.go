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
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

func TestJenkinsStagesDataFlow(t *testing.T) {
	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake`,
		},
	}

	// import raw data table
	// SELECT * FROM _raw_jenkins_api_stages INTO OUTFILE "/tmp/_raw_jenkins_api_stages.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_stages.csv", "_raw_jenkins_api_stages")

	// import tool table
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jenkins_builds.csv", models.JenkinsBuild{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jenkins_build_commits.csv", models.JenkinsBuildCommit{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/cicd_tasks.csv", devops.CICDTask{})

	// verify extraction
	dataflowTester.FlushTabler(&models.JenkinsStage{})
	dataflowTester.Subtask(tasks.ExtractApiStagesMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsStage{},
		"./snapshot_tables/_tool_jenkins_stages.csv",
		[]string{
			"connection_id",
			"id",
			"name",
			"exec_node",
			"status",
			"start_time_millis",
			"duration_millis",
			"pause_duration_millis",
			"build_name",
			"type",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.Subtask(tasks.ConvertStagesMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks_after_stages.csv",
		[]string{
			"id",
			"name",
			"pipeline_id",
			"result",
			"status",
			"type",
			"duration_sec",
			"started_date",
			"finished_date",
			"environment",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
