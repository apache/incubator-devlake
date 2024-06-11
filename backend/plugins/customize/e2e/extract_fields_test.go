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

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/customize/service"
	"github.com/apache/incubator-devlake/plugins/customize/tasks"
)

func TestExtractFieldDataExplicitFlow1(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	taskData := &tasks.TaskData{
		Options: &tasks.Options{
			TransformationRules: []tasks.MappingRules{{
				Table:         "issues",
				RawDataTable:  "_raw_jira_api_issues",
				RawDataParams: `{"ConnectionId":1,"BoardId":8}`,
				Mapping: map[string]string{
					"x_test":  "fields.created",
					"x_float": "fields.customfield_10024",
					"x_int":   "fields.customfield_10146",
				},
			}}}}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issues.csv", "_raw_jira_api_issues")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issues.csv", &ticket.Issue{})
	dataflowTester.FlushTabler(&models.CustomizedField{})
	svc := service.NewService(dataflowTester.Dal)
	err := svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_test",
		DisplayName: "test column",
		DataType:    "varchar(255)",
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
		ColumnName:  "x_int",
		DisplayName: "test column x_int",
		DataType:    "bigint",
	})
	if err != nil {
		t.Fatal(err)
	}
	// verify extension fields extraction
	dataflowTester.Subtask(tasks.ExtractCustomizedFieldsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		"./snapshot_tables/issues_for_data_extraction_board8.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"x_test",
			"x_float",
			"x_int",
		),
	)
}

func TestExtractFieldDataExplicitFlow2(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	taskData := &tasks.TaskData{
		Options: &tasks.Options{
			TransformationRules: []tasks.MappingRules{{
				Table:         "issues",
				RawDataTable:  "_raw_jira_api_issues",
				RawDataParams: `{"ConnectionId":1,"BoardId":9}`,
				Mapping: map[string]string{
					"x_test":  "fields.created",
					"x_float": "fields.customfield_10024",
					"x_int":   "fields.customfield_10146",
				},
			}}}}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issues.csv", "_raw_jira_api_issues")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issues.csv", &ticket.Issue{})
	dataflowTester.FlushTabler(&models.CustomizedField{})
	svc := service.NewService(dataflowTester.Dal)
	err := svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_test",
		DisplayName: "test column",
		DataType:    "varchar(255)",
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
		ColumnName:  "x_int",
		DisplayName: "test column x_int",
		DataType:    "bigint",
	})
	if err != nil {
		t.Fatal(err)
	}
	// verify extension fields extraction
	dataflowTester.Subtask(tasks.ExtractCustomizedFieldsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		"./snapshot_tables/issues_for_data_extraction_board9.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"x_test",
			"x_float",
			"x_int",
		),
	)
}

func TestExtractFieldDataWildcardFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	taskData := &tasks.TaskData{
		Options: &tasks.Options{
			TransformationRules: []tasks.MappingRules{{
				Table:         "issues",
				RawDataTable:  "_raw_jira_api_issues",
				RawDataParams: `{"ConnectionId":1,%}`,
				Mapping: map[string]string{
					"x_test":  "fields.created",
					"x_float": "fields.customfield_10024",
					"x_int":   "fields.customfield_10146",
				},
			}}}}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issues.csv", "_raw_jira_api_issues")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issues.csv", &ticket.Issue{})
	dataflowTester.FlushTabler(&models.CustomizedField{})
	svc := service.NewService(dataflowTester.Dal)
	err := svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_test",
		DisplayName: "test column",
		DataType:    "varchar(255)",
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
		ColumnName:  "x_int",
		DisplayName: "test column x_int",
		DataType:    "bigint",
	})
	if err != nil {
		t.Fatal(err)
	}
	// verify extension fields extraction
	dataflowTester.Subtask(tasks.ExtractCustomizedFieldsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		"./snapshot_tables/issues_for_data_extraction_wildcard.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"x_test",
			"x_float",
			"x_int",
		),
	)
}
