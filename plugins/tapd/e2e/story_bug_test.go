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
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdStoryBugDataFlow(t *testing.T) {

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
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_story_bugs.csv",
		"_raw_tapd_api_story_bugs")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdStoryBug{})
	dataflowTester.Subtask(tasks.ExtractStoryBugsMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdStoryBug{},
		"./snapshot_tables/_tool_tapd_story_bugs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"workspace_id",
			"story_id",
			"bug_id",
			"issue_updated",
		),
	)

}
