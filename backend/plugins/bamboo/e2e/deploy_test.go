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

func TestBambooDeployDataFlow(t *testing.T) {
	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)
	taskData := &tasks.BambooTaskData{
		Options: &models.BambooOptions{
			ConnectionId:      1,
			PlanKey:           "TEST1",
			BambooScopeConfig: new(models.BambooScopeConfig),
		},
		ApiClient: getFakeAPIClient(),
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bamboo_api_deploys.csv", "_raw_bamboo_api_deploys")
	// it needs import plan data
	//dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_bamboo_plans_for_deploys.csv", models.BambooPlan{})

	// verify extraction
	dataflowTester.FlushTabler(&models.BambooDeployEnvironment{})
	dataflowTester.Subtask(tasks.ExtractDeployMeta, taskData)
	dataflowTester.VerifyTable(
		models.BambooDeployEnvironment{},
		"./snapshot_tables/_tool_bamboo_deploy_environments.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"env_id",
			"name",
			"plan_key",
			"description",
			"position",
			"configuration_state",
			"can_view",
			"can_edit",
			"can_delete",
			"allowed_to_execute",
			"can_execute",
			"allowed_to_create_version",
			"allowed_to_set_version_status",
		),
	)

}
