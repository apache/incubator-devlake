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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/qa"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestImportQaTestCasesDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// Flush the relevant tables
	dataflowTester.FlushTabler(&qa.QaTestCase{})
	dataflowTester.FlushTabler(&qa.QaProject{}) // qaTestCaseHandler also creates/updates QaProject
	dataflowTester.FlushTabler(&qa.QaApi{})     // qaTestCaseHandler also creates/updates QaApi for API test cases
	dataflowTester.FlushTabler(&qa.QaTestCaseExecution{})
	dataflowTester.FlushTabler(&crossdomain.Account{}) // qaTestCaseHandler also creates/updates Account for API test cases

	// Create a new service instance
	svc := service.NewService(dataflowTester.Dal)

	// Assume a dummy CSV file exists for testing
	// You will need to create this file at backend/plugins/customize/e2e/raw_tables/qa_test_cases_input.csv
	// with appropriate test data.
	qaTestCasesFile, err := os.Open("raw_tables/qa_test_cases_input.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer qaTestCasesFile.Close()

	// Define dummy qaProjectId and qaProjectName
	qaProjectId := "test-qa-project-id"
	qaProjectName := "Test QA Project"

	// Import data from the CSV file
	err = svc.ImportQaTestCases(qaProjectId, qaProjectName, qaTestCasesFile, false) // Use false for initial import
	if err != nil {
		t.Fatal(err)
	}

	// Verify the imported data in qa_test_cases table
	// You will need to create the expected snapshot file at backend/plugins/customize/e2e/snapshot_tables/qa_test_cases_output.csv
	// and define the columns to verify based on your qa_test_cases_input.csv and the qa.QaTestCase struct.
	dataflowTester.VerifyTableWithRawData(
		&qa.QaTestCase{},
		"snapshot_tables/qa_test_cases_output.csv",
		[]string{
			"id",
			"name",
			"create_time",
			"creator_id",
			"type",
			"qa_project_id",
			"qa_api_id",
		})

	// Add incremental import test
	qaTestCasesIncrementalFile, err := os.Open("raw_tables/qa_test_cases_input_incremental.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer qaTestCasesIncrementalFile.Close()
	err = svc.ImportQaTestCases(qaProjectId, qaProjectName, qaTestCasesIncrementalFile, true) // Use true for incremental import
	if err != nil {
		t.Fatal(err)
	}
	dataflowTester.VerifyTableWithRawData(
		&qa.QaTestCase{},
		"snapshot_tables/qa_test_cases_output_incremental.csv",
		[]string{
			"id",
			"name",
			"create_time",
			"creator_id",
			"type",
			"qa_project_id",
			"qa_api_id",
		})

	// verify qa_projects and qa_apis
	dataflowTester.VerifyTableWithRawData(
		&qa.QaProject{},
		"snapshot_tables/qa_projects_output.csv",
		[]string{
			"id",
			"name",
		})
	dataflowTester.VerifyTableWithRawData(
		&crossdomain.Account{},
		"snapshot_tables/accounts_from_qa_test_cases.csv",
		[]string{
			"id",
			"full_name",
			"user_name",
		},
	)
}

func TestImportQaTestCasesDataCleanup(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// Flush all relevant tables
	dataflowTester.FlushTabler(&qa.QaTestCase{})
	dataflowTester.FlushTabler(&qa.QaProject{})
	dataflowTester.FlushTabler(&qa.QaApi{})
	dataflowTester.FlushTabler(&qa.QaTestCaseExecution{})
	dataflowTester.FlushTabler(&crossdomain.Account{})

	svc := service.NewService(dataflowTester.Dal)

	qaProjectId := "test-cleanup-project"
	qaProjectName := "Test Cleanup Project"

	// 1. First import import test cases with API references
	testCasesFile, err := os.Open("raw_tables/qa_full_test_cases_input.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer testCasesFile.Close()

	err = svc.ImportQaTestCases(qaProjectId, qaProjectName, testCasesFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// import test case executions
	testCaseExecutionsFile, err := os.Open("raw_tables/qa_test_case_executions_input.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer testCaseExecutionsFile.Close()

	err = svc.ImportQaTestCaseExecutions(qaProjectId, testCaseExecutionsFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// Then import APIs
	apisFile, err := os.Open("raw_tables/qa_apis_input.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer apisFile.Close()

	err = svc.ImportQaApis(qaProjectId, apisFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify APIs, test cases and test case executions were imported
	var initialApiCount int64
	initialApiCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaApi{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}
	if initialApiCount == 0 {
		t.Error("Expected API data to be imported initially")
	}

	var initialTestCaseCount int64
	initialTestCaseCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaTestCase{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}
	if initialTestCaseCount == 0 {
		t.Error("Expected test cases to be imported initially")
	}

	var initialTestCaseExecutionCount int64
	initialTestCaseExecutionCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaTestCaseExecution{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}

	if initialTestCaseExecutionCount == 0 {
		t.Error("Expected test case executions to be imported initially")
	}

	// 2. Second import non-incremental - test cases
	nonApiDataFile, err := os.Open("raw_tables/qa_non_api_test_cases.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer nonApiDataFile.Close()

	err = svc.ImportQaTestCases(qaProjectId, qaProjectName, nonApiDataFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify API data was cleaned up
	var finalApiCount int64
	finalApiCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaApi{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}
	if finalApiCount != 0 {
		t.Errorf("Expected API data to be cleaned up, but found %d records", finalApiCount)
	}

	// Verify test case execution data was cleaned up
	var finalTestCaseExecutionCount int64
	finalTestCaseExecutionCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaTestCaseExecution{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}
	if finalTestCaseExecutionCount != 0 {
		t.Errorf("Expected test case executions to be cleaned up, but found %d records", finalTestCaseExecutionCount)
	}

	// Verify test case count is correct (should be 2)
	var finalTestCaseCount int64
	finalTestCaseCount, err = dataflowTester.Dal.Count(
		dal.From(&qa.QaTestCase{}),
		dal.Where("qa_project_id = ?", qaProjectId),
	)
	if err != nil {
		t.Fatal(err)
	}
	if finalTestCaseCount != 2 {
		t.Errorf("Expected 2 test cases after non-incremental import, got %d", finalTestCaseCount)
	}
}
