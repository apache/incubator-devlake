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

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdWorklogDataFlow(t *testing.T) {

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
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_worklogs.csv",
		"_raw_tapd_api_worklogs")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdWorklog{})
	dataflowTester.Subtask(tasks.ExtractWorklogMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdWorklog{},
		"./snapshot_tables/_tool_tapd_worklogs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"workspace_id",
			"entity_type",
			"entity_id",
			"timespent",
			"spentdate",
			"owner",
			"created",
			"memo",
		),
	)

	dataflowTester.FlushTabler(&ticket.IssueWorklog{})
	dataflowTester.Subtask(tasks.ConvertWorklogMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueWorklog{},
		"./snapshot_tables/issue_worklogs.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"author_id",
			"comment",
			"time_spent_minutes",
			"logged_date",
			"started_date",
			"issue_id",
		),
	)

}
