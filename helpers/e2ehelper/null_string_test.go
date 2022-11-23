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
	"os"
	"path/filepath"
	"testing"
	"time"
)

type nullStringTest struct {
	Id                string    `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	RawDataParams     string    `gorm:"column:_raw_data_params;type:varchar(255);index" json:"_raw_data_params"`
	RawDataTable      string    `gorm:"column:_raw_data_table;type:varchar(255)" json:"_raw_data_table"`
	RawDataId         uint64    `gorm:"column:_raw_data_id" json:"_raw_data_id"`
	RawDataRemark     string    `gorm:"column:_raw_data_remark" json:"_raw_data_remark"`
	IssueId           string    `gorm:"index;type:varchar(255)"`
	AuthorId          string    `gorm:"type:varchar(255)"`
	AuthorName        string    `gorm:"type:varchar(255)"`
	FieldId           string    `gorm:"type:varchar(255)"`
	FieldName         string    `gorm:"type:varchar(255)"`
	OriginalFromValue string
	OriginalToValue   string
	FromValue         string
	ToValue           string
	CreatedDate       time.Time
}

func (nullStringTest) TableName() string {
	return "null_string_tests"
}
func TestNullStringFlow(t *testing.T) {
	//var plugin impl.Jira
	dataflowTester := NewDataFlowTester(t, "", nil)
	dir, err := os.MkdirTemp("", "testnullstring")
	if err != nil {
		panic(err)
	}
	// nolint
	defer os.RemoveAll(dir)

	// verify changelog conversion
	dataflowTester.ImportCsvIntoTabler("./testdata/issue_changelogs.csv", &nullStringTest{})
	dataflowTester.CreateSnapshot(
		&nullStringTest{},
		TableOptions{
			CSVRelPath: filepath.Join(dir, "issue_changelogs_null_string.csv"),
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
