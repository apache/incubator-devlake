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

func TestZentaoTaskWorklogDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    48,
		},
		ApiClient: getFakeAPIClient(),
	}

	// import _raw_zentao_api_task_worklogs raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_task_worklogs.csv",
		"_raw_zentao_api_task_worklogs")
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_zentao_accounts.csv", &models.ZentaoAccount{})

	// verify worklogs extraction
	dataflowTester.FlushTabler(&models.ZentaoWorklog{})
	dataflowTester.Subtask(tasks.ExtractTaskWorklogsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoWorklog{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_worklogs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify task repo commit conversion
	dataflowTester.FlushTabler(&ticket.IssueWorklog{})
	dataflowTester.Subtask(tasks.ConvertTaskWorklogsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.IssueWorklog{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_worklogs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
