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
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/ae/impl"
	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/ae/tasks"
)

func TestAECommitDataFlow(t *testing.T) {
	var ae impl.AE
	dataflowTester := e2ehelper.NewDataFlowTester(t, "ae", ae)

	taskData := &tasks.AeTaskData{
		Options: &tasks.AeOptions{
			ConnectionId: 1,
			ProjectId:    13,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_ae_commits.csv", "_raw_ae_commits")

	// verify extraction
	dataflowTester.FlushTabler(&models.AECommit{})
	dataflowTester.Subtask(tasks.ExtractCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		models.AECommit{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.AECommit{}.TableName()),
		[]string{
			"hex_sha",
			"analysis_id",
			"author_email",
			"dev_eq",
			"ae_project_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&code.Commit{})
	// import data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_raw_commits.csv", &code.Commit{})
	dataflowTester.Subtask(tasks.ConvertCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		code.Commit{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.Commit{}.TableName()),
		[]string{
			"sha",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"additions",
			"deletions",
			"dev_eq",
			"message",
			"author_name",
			"author_email",
			"authored_date",
			"author_id",
			"committer_name",
			"committer_email",
			"committed_date",
			"committer_id",
		},
	)
}
