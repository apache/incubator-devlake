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
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/impl"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/tasks"
)

func TestAzuredevopsBuildDataFlow(t *testing.T) {

	var azuredevops impl.Azuredevops
	dataflowTester := e2ehelper.NewDataFlowTester(t, "azuredevops_go", azuredevops)

	regexEnricher := api.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "(?i)deploy")
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "(?i)prod")

	taskData := &tasks.AzuredevopsTaskData{
		Options: &tasks.AzuredevopsOptions{
			ConnectionId:   1,
			ProjectId:      "test-project",
			OrganizationId: "johndoe",
			RepositoryId:   "0d50ba13-f9ad-49b0-9b21-d29eda50ca33",
			RepositoryType: models.RepositoryTypeADO,
			ScopeConfig:    new(models.AzuredevopsScopeConfig),
		},
		RegexEnricher: regexEnricher,
	}

	dataflowTester.FlushTabler(&models.AzuredevopsBuild{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azuredevops_go_api_builds.csv", "_raw_azuredevops_go_api_builds")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_azuredevops_go_repos.csv", &models.AzuredevopsRepo{})
	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)

	// Omit the datetime columns to avoid test failures caused by the varying precision
	// in MySQL’s datetime(3) and PostgreSQL’s 'timestamp with time zone'.
	dataflowTester.VerifyTable(
		models.AzuredevopsBuild{},
		"./snapshot_tables/_tool_azuredevops_go_builds.csv",
		[]string{
			"connection_id",
			"azuredevops_id",
			"repository_id",
			"status",
			"result",
			"name",
			"source_branch",
			"source_version",
			"tags",
		},
	)

	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})
	dataflowTester.Subtask(tasks.ConvertBuildsMeta, taskData)

	// Omit the datetime columns to avoid test failures caused by the varying precision
	// in MySQL’s datetime(3) and PostgreSQL’s 'timestamp with time zone'.
	dataflowTester.VerifyTable(
		devops.CICDPipeline{},
		"./snapshot_tables/cicd_pipelines.csv",
		[]string{
			"id",
			"name",
			"result",
			"status",
			"original_status",
			"original_result",
			"type",
			"duration_sec",
			"queued_duration_sec",
			"environment",
			"cicd_scope_id",
		},
	)

	dataflowTester.VerifyTableWithOptions(&devops.CiCDPipelineCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipeline_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

}
