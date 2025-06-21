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
	"os"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestImportIssueDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// import raw data table
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&models.CustomizedField{})
	dataflowTester.FlushTabler(&ticket.IssueLabel{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.FlushTabler(&ticket.SprintIssue{})
	svc := service.NewService(dataflowTester.Dal)
	err := svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_varchar",
		DisplayName: "test column x_varchar",
		DataType:    "varchar(255)",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_text",
		DisplayName: "test column x_text",
		DataType:    "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_int",
		DisplayName: "test column x_int",
		DataType:    "bigint",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_float",
		DisplayName: "test column x_float",
		DataType:    "float",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_time",
		DisplayName: "test column x_time",
		DataType:    "timestamp",
	})
	if err != nil {
		t.Fatal(err)
	}

	issueFile, err1 := os.Open("raw_tables/issues_input.csv")
	if err1 != nil {
		t.Fatal(err1)
	}
	defer issueFile.Close()
	err = svc.ImportIssue("csv-board", issueFile, false)
	if err != nil {
		t.Fatal(err)
	}
	issueAppendToFile1, err2 := os.Open("raw_tables/issues_input_incremental.csv")
	if err2 != nil {
		t.Fatal(err2)
	}
	defer issueAppendToFile1.Close()
	err = svc.ImportIssue("csv-board", issueAppendToFile1, true)
	if err != nil {
		t.Fatal(err)
	}
	issueFile2, err3 := os.Open("raw_tables/issues_input2.csv")
	if err3 != nil {
		t.Fatal(err3)
	}
	defer issueFile2.Close()
	err = svc.ImportIssue("csv-board2", issueFile2, false)
	if err != nil {
		t.Fatal(err)
	}
	issueToOverwriteFile2, err4 := os.Open("raw_tables/issues_input2_overwrite.csv")
	if err4 != nil {
		t.Fatal(err4)
	}
	defer issueToOverwriteFile2.Close()
	err = svc.ImportIssue("csv-board2", issueToOverwriteFile2, false)
	if err != nil {
		t.Fatal(err)
	}
	dataflowTester.VerifyTableWithRawData(
		ticket.Issue{},
		"snapshot_tables/issues_output.csv",
		[]string{
			"id",
			"url",
			"icon_url",
			"issue_key",
			"title",
			"description",
			"epic_key",
			"type",
			"original_type",
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
			"creator_name",
			"assignee_id",
			"assignee_name",
			"severity",
			"component",
			"original_project",
			"x_varchar",
			"x_text",
			"x_time",
			"x_float",
			"x_int",
		})
	dataflowTester.VerifyTableWithRawData(
		&ticket.IssueLabel{},
		"snapshot_tables/issue_labels.csv",
		[]string{
			"issue_id",
			"label_name",
		})
	dataflowTester.VerifyTableWithRawData(
		&ticket.BoardIssue{},
		"snapshot_tables/board_issues.csv",
		[]string{
			"board_id",
			"issue_id",
		})
	dataflowTester.VerifyTableWithRawData(
		&crossdomain.Account{},
		"snapshot_tables/accounts.csv",
		[]string{
			"id",
			"full_name",
			"user_name",
		},
	)

	dataflowTester.VerifyTableWithRawData(
		&ticket.SprintIssue{},
		"snapshot_tables/sprint_issues.csv",
		[]string{
			"sprint_id",
			"issue_id",
		},
	)
}
