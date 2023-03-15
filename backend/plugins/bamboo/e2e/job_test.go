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
	"github.com/apache/incubator-devlake/plugins/bamboo/impl"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

func TestBambooJobDataFlow(t *testing.T) {

	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)

	taskData := &tasks.BambooTaskData{
		Options: &models.BambooOptions{
			ConnectionId: 3,
			ProjectKey:   "TEST1",
			BambooTransformationRule: &models.BambooTransformationRule{
				DeploymentPattern: "(?i)compile",
				ProductionPattern: "(?i)compile",
			},
		},
	}
	// import raw data table
	// SELECT * FROM _raw_bamboo_api_job INTO OUTFILE "/tmp/_raw_bamboo_api_job.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bamboo_api_job.csv", "_raw_bamboo_api_job")

	// verify extraction
	dataflowTester.FlushTabler(&models.BambooJob{})
	dataflowTester.Subtask(tasks.ExtractJobMeta, taskData)
	dataflowTester.VerifyTable(
		models.BambooJob{},
		"./snapshot_tables/_tool_bamboo_jobs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"job_key",
			"id",
			"name",
			"plan_key",
			"project_key",
			"project_name",
			"description",
			"branch_name",
			"stage_name",
			"type",
		),
	)

	// verify conversion
	/*dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.Subtask(tasks.ConvertJobMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDTask{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_tasks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})*/
}
