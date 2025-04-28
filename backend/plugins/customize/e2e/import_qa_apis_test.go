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

func TestImportQaApisDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// Flush the qa_apis table
	dataflowTester.FlushTabler(&qa.QaApi{})
	dataflowTester.FlushTabler(&crossdomain.Account{})

	// Create a new service instance
	svc := service.NewService(dataflowTester.Dal)

	// Assume a dummy CSV file exists for testing
	// You will need to create this file at backend/plugins/customize/e2e/raw_tables/qa_apis_input.csv
	// with appropriate test data.
	qaApisFile, err := os.Open("raw_tables/qa_apis_input.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer qaApisFile.Close()

	// Define a dummy qaProjectId
	qaProjectId := "test-qa-project-id"

	// Import data from the CSV file
	err = svc.ImportQaApis(qaProjectId, qaApisFile, false) // Use false for initial import
	if err != nil {
		t.Fatal(err)
	}

	// Verify the imported data
	// You will need to create the expected snapshot file at backend/plugins/customize/e2e/snapshot_tables/qa_apis_output.csv
	// and define the columns to verify based on your qa_apis_input.csv and the qa.QaApi struct.
	dataflowTester.VerifyTableWithRawData(
		&qa.QaApi{},
		"snapshot_tables/qa_apis_output.csv",
		[]string{
			"id",
			"name",
			"path",
			"method",
			"create_time",
			"creator_id",
			"qa_project_id",
		})

	dataflowTester.VerifyTableWithRawData(
		&crossdomain.Account{},
		"snapshot_tables/accounts_from_qa_apis_output.csv",
		[]string{
			"id",
			"full_name",
			"user_name",
		})
	// Add incremental import test if needed, similar to import_issues_test.go
	qaApisIncrementalFile, err := os.Open("raw_tables/qa_apis_input_incremental.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer qaApisIncrementalFile.Close()
	err = svc.ImportQaApis(qaProjectId, qaApisIncrementalFile, true) // Use true for incremental import
	if err != nil {
		t.Fatal(err)
	}
	dataflowTester.VerifyTableWithRawData(
		&qa.QaApi{},
		"snapshot_tables/qa_apis_output_incremental.csv",
		[]string{
			"id",
			"name",
			"path",
			"method",
			"create_time",
			"creator_id",
			"qa_project_id",
		})
	dataflowTester.VerifyTableWithRawData(
		&crossdomain.Account{},
		"snapshot_tables/accounts_from_qa_apis_incremental_output.csv",
		[]string{
			"id",
			"full_name",
			"user_name",
		},
	)
}
