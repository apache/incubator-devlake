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

func TestTapdStoryDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_story_status.csv",
		"_raw_tapd_api_story_status")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdStoryStatus{})
	dataflowTester.Subtask(tasks.ExtractStoryStatusMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdStoryStatus{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdStoryStatus{}.TableName()),
		[]string{"connection_id", "workspace_id", "english_name"},
		[]string{
			"chinese_name",
			"is_last_step",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_stories.csv",
		"_raw_tapd_api_stories")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdStory{})
	dataflowTester.FlushTabler(&models.TapdWorkSpaceStory{})
	dataflowTester.FlushTabler(&models.TapdIterationStory{})
	dataflowTester.FlushTabler(&models.TapdStoryLabel{})
	dataflowTester.Subtask(tasks.ExtractStoryMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdStory{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdStory{}.TableName()),
		[]string{"connection_id", "id"},
		[]string{
			"workitem_type_id",
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
			"size",
			"priority",
			"developer",
			"iteration_id",
			"test_focus",
			"type",
			"source",
			"module",
			"version",
			"completed",
			"category_id",
			"path",
			"parent_id",
			"children_id",
			"ancestor_id",
			"business_value",
			"effort",
			"effort_completed",
			"exceed",
			"remain",
			"release_id",
			"confidential",
			"templated_id",
			"created_from",
			"feature",
			"std_status",
			"std_type",
			"url",
			"attachment_count",
			"has_attachment",
			"bug_id",
			"follower",
			"sync_type",
			"predecessor_count",
			"is_archived",
			"modifier",
			"progress_manual",
			"successor_count",
			"label",
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
		models.TapdWorkSpaceStory{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdWorkSpaceStory{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"story_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdIterationStory{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdIterationStory{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"iteration_id",
			"story_id",
		},
		[]string{
			"resolution_date",
			"story_created_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdStoryLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdStoryLabel{}.TableName()),
		[]string{
			"label_name",
			"story_id",
		},
		[]string{
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
	dataflowTester.Subtask(tasks.ConvertStoryMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		fmt.Sprintf("./snapshot_tables/%s_story.csv", ticket.Issue{}.TableName()),
		[]string{"id"},
		[]string{
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
		fmt.Sprintf("./snapshot_tables/%s_story.csv", ticket.BoardIssue{}.TableName()),
		[]string{
			"board_id",
			"issue_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		ticket.SprintIssue{},
		fmt.Sprintf("./snapshot_tables/%s_story.csv", ticket.SprintIssue{}.TableName()),
		[]string{
			"issue_id",
			"sprint_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.Subtask(tasks.ConvertStoryLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueLabel{},
		fmt.Sprintf("./snapshot_tables/%s_story.csv", ticket.IssueLabel{}.TableName()),
		[]string{
			"issue_id",
			"label_name",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

}
