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

func TestJenkinsJobsDataFlow(t *testing.T) {

	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_jobs.csv", "_raw_jenkins_api_jobs")

	// verify extraction
	dataflowTester.FlushTabler(&models.JenkinsJob{})
	dataflowTester.Subtask(tasks.ExtractApiJobsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsJob{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.JenkinsJob{}.TableName()),
		[]string{
			"connection_id",
			"name",
			"path",
			"class",
			"color",
			"base",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&devops.Job{})
	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTable(
		devops.Job{},
		fmt.Sprintf("./snapshot_tables/%s.csv", devops.Job{}.TableName()),
		[]string{
			"name",
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
