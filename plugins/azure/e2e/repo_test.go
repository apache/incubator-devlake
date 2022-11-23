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

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/azure/impl"
	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/azure/tasks"
)

func TestAzureRepoDataFlow(t *testing.T) {
	var github impl.Azure
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", github)

	taskData := &tasks.AzureTaskData{
		Options: &tasks.AzureOptions{
			ConnectionId: 1,
			Project:      "test",
		},
	}

	// import raw data table
	// SELECT * FROM _raw_azure_api_repositories INTO OUTFILE "/tmp/_raw_azure_api_repositories.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azure_api_repositories.csv", "_raw_azure_api_repositories")

	// verify extraction
	dataflowTester.FlushTabler(&models.AzureRepo{})
	dataflowTester.Subtask(tasks.ExtractApiRepoMeta, taskData)
	dataflowTester.VerifyTable(
		models.AzureRepo{},
		"./snapshot_tables/_tool_azure_repos.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"azure_id",
			"name",
			"url",
			"project_id",
			"default_branch",
			"size",
			"remote_url",
			"ssh_url",
			"web_url",
			"is_disabled",
		),
	)
}
