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
	"github.com/apache/incubator-devlake/plugins/ae/impl"
	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/ae/tasks"
)

func TestAEProjectDataFlow(t *testing.T) {
	var ae impl.AE
	dataflowTester := e2ehelper.NewDataFlowTester(t, "ae", ae)

	taskData := &tasks.AeTaskData{
		Options: &tasks.AeOptions{
			ConnectionId: 1,
			ProjectId:    13,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_ae_project.csv", "_raw_ae_project")

	// verify extraction
	dataflowTester.FlushTabler(&models.AEProject{})
	dataflowTester.Subtask(tasks.ExtractProjectMeta, taskData)
	dataflowTester.VerifyTable(
		models.AEProject{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.AEProject{}.TableName()),
		[]string{"connection_id", "id"},
		[]string{
			"git_url",
			"priority",
			"ae_create_time",
			"ae_update_time",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
