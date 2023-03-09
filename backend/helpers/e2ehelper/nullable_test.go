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

package e2ehelper

import (
	"testing"
)

func TestNullableFlow(t *testing.T) {
	dataflowTester := NewDataFlowTester(t, "", nil)

	// nullable as false
	dataflowTester.ImportCsvIntoTabler("./testdata/issue_changelogs.csv", &nullStringTest{})
	dataflowTester.VerifyTableWithOptions(
		&nullStringTest{},
		TableOptions{
			CSVRelPath: "./testdata/issue_changelogs.csv",
			TargetFields: ColumnWithRawData(
				"id",
				"issue_id",
				"author_id",
				"author_name",
				"field_id",
				"field_name",
				"original_from_value",
				"original_to_value",
				"from_value",
				"to_value",
				"created_date",
				"_raw_data_params",
				"_raw_data_table",
				"_raw_data_id",
				"_raw_data_remark",
			),
		},
	)
	// nullable as true
	dataflowTester.ImportNullableCsvIntoTabler("./testdata/issue_changelogs.csv", &nullStringTest{})
	dataflowTester.VerifyTableWithOptions(
		&nullStringTest{},
		TableOptions{
			Nullable:   true,
			CSVRelPath: "./testdata/issue_changelogs.csv",
			TargetFields: ColumnWithRawData(
				"id",
				"issue_id",
				"author_id",
				"author_name",
				"field_id",
				"field_name",
				"original_from_value",
				"original_to_value",
				"from_value",
				"to_value",
				"created_date",
				"_raw_data_params",
				"_raw_data_table",
				"_raw_data_id",
				"_raw_data_remark",
			),
		},
	)
}
