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

func TestTapdBugDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}

	// bug status
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_bug_status.csv",
		"_raw_tapd_api_bug_status")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdBugStatus{})
	dataflowTester.Subtask(tasks.ExtractBugStatusMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdBugStatus{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdBugStatus{}.TableName()),
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
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_bugs.csv",
		"_raw_tapd_api_bugs")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdBug{})
	dataflowTester.FlushTabler(&models.TapdWorkSpaceBug{})
	dataflowTester.FlushTabler(&models.TapdIterationBug{})
	dataflowTester.FlushTabler(&models.TapdBugLabel{})
	dataflowTester.Subtask(tasks.ExtractBugMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdBug{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdBug{}.TableName()),
		[]string{"connection_id", "id"},
		[]string{
			"epic_key",
			"title",
			"description",
			"workspace_id",
			"created",
			"modified",
			"status",
			"cc",
			"begin",
			"due",
			"priority",
			"iteration_id",
			"source",
			"module",
			"release_id",
			"created_from",
			"feature",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"severity",
			"reporter",
			"resolved",
			"closed",
			"lastmodify",
			"auditer",
			"de",
			"fixer",
			"version_test",
			"version_report",
			"version_close",
			"version_fix",
			"baseline_find",
			"baseline_join",
			"baseline_close",
			"baseline_test",
			"sourcephase",
			"te",
			"current_owner",
			"resolution",
			"originphase",
			"confirmer",
			"participator",
			"closer",
			"platform",
			"os",
			"testtype",
			"testphase",
			"frequency",
			"regression_number",
			"flows",
			"testmode",
			"issue_id",
			"verify_time",
			"reject_time",
			"reopen_time",
			"audit_time",
			"suspend_time",
			"deadline",
			"in_progress_time",
			"assigned_time",
			"template_id",
			"story_id",
			"std_status",
			"std_type",
			"type",
			"url",
			"support_id",
			"support_forum_id",
			"ticket_id",
			"follower",
			"sync_type",
			"label",
			"effort",
			"effort_completed",
			"exceed",
			"remain",
			"progress",
			"estimate",
			"bugtype",
			"milestone",
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
		},
	)
	dataflowTester.VerifyTable(
		models.TapdWorkSpaceBug{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdWorkSpaceBug{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"bug_id",
		},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdIterationBug{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdIterationBug{}.TableName()),
		[]string{
			"connection_id",
			"workspace_id",
			"iteration_id",
			"bug_id",
		},
		[]string{
			"resolution_date",
			"bug_created_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.TapdBugLabel{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.TapdBugLabel{}.TableName()),
		[]string{
			"label_name",
			"bug_id",
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
	dataflowTester.Subtask(tasks.ConvertBugMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		fmt.Sprintf("./snapshot_tables/%s_bug.csv", ticket.Issue{}.TableName()),
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
		fmt.Sprintf("./snapshot_tables/%s_bug.csv", ticket.BoardIssue{}.TableName()),
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
		fmt.Sprintf("./snapshot_tables/%s_bug.csv", ticket.SprintIssue{}.TableName()),
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
	dataflowTester.Subtask(tasks.ConvertBugLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueLabel{},
		fmt.Sprintf("./snapshot_tables/%s_bug.csv", ticket.IssueLabel{}.TableName()),
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
