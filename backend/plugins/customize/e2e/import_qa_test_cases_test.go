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
	dataflowTester.FlushTabler(&qa.QaProject{})        // qaTestCaseHandler also creates/updates QaProject
	dataflowTester.FlushTabler(&qa.QaApi{})            // qaTestCaseHandler also creates/updates QaApi for API test cases
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
