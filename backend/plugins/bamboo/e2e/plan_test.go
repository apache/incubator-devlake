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

func TestBambooPlanDataFlow(t *testing.T) {

	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)

	taskData := &tasks.BambooTaskData{
		Options: &models.BambooOptions{
			ConnectionId:             3,
			ProjectKey:               "TEST1",
			BambooTransformationRule: new(models.BambooTransformationRule),
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bamboo_api_plan.csv", "_raw_bamboo_api_plan")

	// verify extraction
	dataflowTester.FlushTabler(&models.BambooPlan{})
	dataflowTester.Subtask(tasks.ExtractPlanMeta, taskData)
	dataflowTester.VerifyTable(
		models.BambooPlan{},
		"./snapshot_tables/_tool_bamboo_plans.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"plan_key",
			"name",
			"expand",
			"project_key",
			"project_name",
			"description",
			"short_name",
			"build_name",
			"short_key",
			"type",
			"enabled",
			"href",
			"rel",
			"is_favourite",
			"is_active",
			"is_building",
			"average_build_time_in_seconds",
		),
	)

}
