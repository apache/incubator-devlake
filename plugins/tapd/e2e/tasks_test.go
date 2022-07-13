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

func TestTapdTaskDataFlow(t *testing.T) {

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
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_tasks.csv",
		"_raw_tapd_api_tasks")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdTask{})
	dataflowTester.FlushTabler(&models.TapdWorkSpaceTask{})
	dataflowTester.FlushTabler(&models.TapdIterationTask{})
	dataflowTester.FlushTabler(&models.TapdTaskLabel{})
	dataflowTester.Subtask(tasks.ExtractTaskMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdTask{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdTask{}.TableName()),
		[]string{
			"connection_id",
			"id",
			"name",
			"description",
			"workspace_id",
			"creator",
			"created",
			"modified",
			"status",
			"owner",
			"cc",
			"begin",
			"due",
			"priority",
			"iteration_id",
			"completed",
			"effort",
			"effort_completed",
			"exceed",
			"remain",
			"std_status",
			"std_type",
			"type",
			"story_id",
			"progress",
			"has_attachment",
			"url",
			"attachment_count",
			"follower",
			"created_from",
			"predecessor_count",
			"successor_count",
			"release_id",
			"label",
			"new_story_id",
			"custom_field_one",
			"custom_field_two",
			"custom_field_three",
			"custom_field_four",
			"custom_field_five",
			"custom_field6",
			"custom_field7",
			"custom_field8",
			"custom_field9",
			"custom_field10",
			"custom_field11",
			"custom_field12",
			"custom_field13",
			"custom_field14",
			"custom_field15",
			"custom_field16",
			"custom_field17",
			"custom_field18",
			"custom_field19",
			"custom_field20",
			"custom_field21",
			"custom_field22",
			"custom_field23",
			"custom_field24",
			"custom_field25",
			"custom_field26",
			"custom_field27",
			"custom_field28",
			"custom_field29",
			"custom_field30",
			"custom_field31",
			"custom_field32",
			"custom_field33",
			"custom_field34",
			"custom_field35",
			"custom_field36",
			"custom_field37",
			"custom_field38",
			"custom_field39",
			"custom_field40",
			"custom_field41",
			"custom_field42",
			"custom_field43",
			"custom_field44",
			"custom_field45",
			"custom_field46",
			"custom_field47",
			"custom_field48",
			"custom_field49",
			"custom_field50",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdWorkSpaceTask{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdWorkSpaceTask{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"task_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdIterationTask{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdIterationTask{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"iteration_id",
			"task_id",
			"resolution_date",
			"task_created_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdTaskLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdTaskLabel{}.TableName()),
		[]string{
			"label_name",
			"task_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.FlushTabler(&ticket.SprintIssue{})
	dataflowTester.FlushTabler(&ticket.IssueLabel{})
	dataflowTester.Subtask(tasks.ConvertTaskMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		fmt.Sprintf("./snapshot_tables/%s_task.csv", ticket.Issue{}.TableName()),
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"url",
			"issue_key",
			"title",
			"description",
			"epic_key",
			"type",
			"status",
			"original_status",
			"story_point",
			"resolution_date",
			"created_date",
			"updated_date",
			"lead_time_minutes",
			"parent_issue_id",
			"priority",
			"original_estimate_minutes",
			"time_spent_minutes",
			"time_remaining_minutes",
			"creator_id",
			"assignee_id",
			"assignee_name",
			"severity",
			"component",
			"icon_url",
			"creator_name",
		},
	)
	dataflowTester.VerifyTable(
		ticket.BoardIssue{},
		fmt.Sprintf("./snapshot_tables/%s_task.csv", ticket.BoardIssue{}.TableName()),
		[]string{
			"board_id",
			"issue_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		ticket.SprintIssue{},
		fmt.Sprintf("./snapshot_tables/%s_task.csv", ticket.SprintIssue{}.TableName()),
		[]string{
			"issue_id",
			"sprint_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.Subtask(tasks.ConvertTaskLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueLabel{},
		fmt.Sprintf("./snapshot_tables/%s_task.csv", ticket.IssueLabel{}.TableName()),
		[]string{
			"issue_id",
			"label_name",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

}
