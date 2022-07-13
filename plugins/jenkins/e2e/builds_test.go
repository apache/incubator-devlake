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
	"testing"

	"github.com/apache/incubator-devlake/models/domainlayer/devops"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

func TestJenkinsBuildsDataFlow(t *testing.T) {

	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_builds.csv", "_raw_jenkins_api_builds")

	// verify extraction
	dataflowTester.FlushTabler(&models.JenkinsBuild{})
	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.JenkinsBuild{}.TableName()),
		[]string{
			"connection_id",
			"job_name",
			"number",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"connection_id",
			"job_name",
			"duration",
			"display_name",
			"estimated_duration",
			"number",
			"result",
			"timestamp",
			"start_time",
			"commit_sha",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&devops.Build{})
	dataflowTester.Subtask(tasks.ConvertBuildsMeta, taskData)
	dataflowTester.VerifyTable(
		devops.Build{},
		fmt.Sprintf("./snapshot_tables/%s.csv", devops.Build{}.TableName()),
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"job_id",
			"name",
			"commit_sha",
			"duration_sec",
			"status",
			"started_date",
		},
	)
}
