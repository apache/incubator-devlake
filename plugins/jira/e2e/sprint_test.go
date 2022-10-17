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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
)

func TestSprintDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_sprints.csv", "_raw_jira_api_sprints")

	// verify sprint extraction
	dataflowTester.FlushTabler(&models.JiraSprint{})
	dataflowTester.FlushTabler(&models.JiraBoardSprint{})
	dataflowTester.Subtask(tasks.ExtractSprintsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraSprint{},
		"./snapshot_tables/_tool_jira_sprints.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"sprint_id",
			"self",
			"state",
			"name",
			"start_date",
			"end_date",
			"complete_date",
			"origin_board_id",
		),
	)

	dataflowTester.VerifyTable(
		models.JiraBoardSprint{},
		"./snapshot_tables/_tool_jira_board_sprints.csv",
		[]string{"connection_id", "board_id", "sprint_id"},
	)

	// verify sprint conversion
	dataflowTester.FlushTabler(&ticket.Sprint{})
	dataflowTester.FlushTabler(&ticket.BoardSprint{})
	dataflowTester.Subtask(tasks.ConvertSprintsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Sprint{},
		"./snapshot_tables/sprints.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"url",
			"status",
			"name",
			"started_date",
			"ended_date",
			"completed_date",
			"original_board_id",
		),
	)
	dataflowTester.VerifyTable(
		ticket.BoardSprint{},
		"./snapshot_tables/board_sprints.csv",
		e2ehelper.ColumnWithRawData(
			"board_id",
			"sprint_id",
		),
	)
}
