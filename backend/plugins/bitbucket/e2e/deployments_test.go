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
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestDeloymentsDataFlow(t *testing.T) {
	var bitbucket impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", bitbucket)

	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			FullName:     "likyh/likyhphp",
			BitbucketTransformationRule: &models.BitbucketTransformationRule{
				DeploymentPattern: "",
				ProductionPattern: "",
			},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_api_deployments.csv", "_raw_bitbucket_api_deployments")

	dataflowTester.FlushTabler(&models.BitbucketDeployment{})
	// verify extraction
	dataflowTester.Subtask(tasks.ExtractApiDeploymentsMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketDeployment{},
		"./snapshot_tables/_tool_bitbucket_deployments.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"bitbucket_id",
			"pipeline_id",
			"step_id",
			"type",
			"name",
			"environment",
			"environment_type",
			"web_url",
			"status",
			"state_url",
			"commit_sha",
			"commit_url",
			"created_on",
			"started_on",
			"completed_on",
			"last_update_time",
			//"key",
		),
	)
}
