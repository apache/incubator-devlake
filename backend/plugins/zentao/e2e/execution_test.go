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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

func TestZentaoExecutionDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    192,
		},
		ApiClient:   getFakeAPIClient(),
		HomePageURL: getFakeHomepage(),
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_executions.csv",
		"_raw_zentao_api_executions")

	// verify extraction
	dataflowTester.FlushTabler(&models.ZentaoExecution{})
	dataflowTester.FlushTabler(&models.ZentaoProductSummary{})
	dataflowTester.Subtask(tasks.ExtractExecutionMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoExecution{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_executions.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&ticket.Sprint{})
	dataflowTester.FlushTabler(&ticket.BoardSprint{})
	dataflowTester.Subtask(tasks.ConvertExecutionMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.Sprint{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/execution_sprint.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.VerifyTableWithOptions(&ticket.BoardSprint{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/execution_board_sprint.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
