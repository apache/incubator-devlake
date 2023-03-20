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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/trello/impl"
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"github.com/apache/incubator-devlake/plugins/trello/tasks"
	"testing"
)

func TestTrelloMemberDataFlow(t *testing.T) {
	var trello impl.Trello
	dataflowTester := e2ehelper.NewDataFlowTester(t, "trello", trello)

	taskData := &tasks.TrelloTaskData{
		Options: &tasks.TrelloOptions{
			ConnectionId: 1,
			BoardId:      "6402f643d23aa9af56b28f4b",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_trello_members.csv", "_raw_trello_members")

	// verify extraction
	dataflowTester.FlushTabler(&models.TrelloMember{})
	dataflowTester.Subtask(tasks.ExtractMemberMeta, taskData)
	dataflowTester.VerifyTableWithOptions(models.TrelloMember{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_trello_members.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
