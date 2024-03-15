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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/impl"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/tasks"
)

func TestAzureDevopsRepoDataFlow(t *testing.T) {

	var azuredevops impl.Azuredevops
	dataflowTester := e2ehelper.NewDataFlowTester(t, "azuredevops_go", azuredevops)

	taskData := &tasks.AzuredevopsTaskData{
		Options: &tasks.AzuredevopsOptions{
			ConnectionId: 1,
			RepositoryId: "0d50ba13-f9ad-49b0-9b21-d29eda50ca33",
			ScopeConfig:  new(models.AzuredevopsScopeConfig),
		},
	}

	// verify conversion
	dataflowTester.FlushTabler(&code.Repo{})
	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_azuredevops_go_repos.csv", &models.AzuredevopsRepo{})
	dataflowTester.Subtask(tasks.ConvertRepoMeta, taskData)
	dataflowTester.VerifyTable(
		code.Repo{},
		"./snapshot_tables/repos.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"description",
			"owner_id",
			"name",
			"url",
			"deleted",
		),
	)

	dataflowTester.VerifyTableWithOptions(&devops.CicdScope{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_scopes.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
