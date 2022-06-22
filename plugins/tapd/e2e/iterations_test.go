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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdIterationDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_iterations.csv",
		"_raw_tapd_api_iterations")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdIteration{})
	dataflowTester.FlushTabler(&models.TapdWorkspaceIteration{})
	dataflowTester.Subtask(tasks.ExtractIterationMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdIteration{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdIteration{}.TableName()),
		[]string{"connection_id", "id"},
		[]string{
			"name",
			"workspace_id",
			"startdate",
			"enddate",
			"status",
			"release_id",
			"description",
			"creator",
			"created",
			"modified",
			"completed",
			"releaseowner",
			"launchdate",
			"notice",
			"releasename",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdWorkspaceIteration{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdWorkspaceIteration{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"iteration_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.FlushTabler(&ticket.Sprint{})
	dataflowTester.FlushTabler(&ticket.BoardSprint{})
	dataflowTester.Subtask(tasks.ConvertIterationMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Sprint{},
		fmt.Sprintf("./snapshot_tables/%s.csv", ticket.Sprint{}.TableName()),
		[]string{"id"},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"name",
			"url",
			"status",
			"started_date",
			"ended_date",
			"completed_date",
			"original_board_id",
		},
	)
	dataflowTester.VerifyTable(
		ticket.BoardSprint{},
		fmt.Sprintf("./snapshot_tables/%s.csv", ticket.BoardSprint{}.TableName()),
		[]string{
			"board_id",
			"sprint_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

}
