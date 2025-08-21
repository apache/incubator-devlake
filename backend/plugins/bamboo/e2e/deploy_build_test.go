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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/impl"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

func getFakeAPIClient() *helper.ApiAsyncClient {
	client := &helper.ApiClient{}
	client.Setup("http://127.0.0.1:8080/bamboo/", nil, time.Second)
	return &helper.ApiAsyncClient{
		ApiClient: client,
	}
}

func TestBambooDeployBuildDataFlow(t *testing.T) {
	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)
	taskData := &tasks.BambooOptions{
		Options: &models.BambooOptions{
			ConnectionId: 1,
			PlanKey:      "TEST-PLA2",
			BambooScopeConfig: &models.BambooScopeConfig{
				DeploymentPattern: "(?i)release",
				ProductionPattern: "(?i)release",
			},
		},
		RegexEnricher: helper.NewRegexEnricher(),
		ApiClient:     getFakeAPIClient(),
	}
	taskData.RegexEnricher.TryAdd(devops.DEPLOYMENT, taskData.Options.DeploymentPattern)
	taskData.RegexEnricher.TryAdd(devops.PRODUCTION, taskData.Options.ProductionPattern)
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bamboo_api_deploy_builds.csv", "_raw_bamboo_api_deploy_builds")

	// verify extraction
	dataflowTester.FlushTabler(&models.BambooDeployBuild{})
	dataflowTester.FlushTabler(&models.BambooPlanBuildVcsRevision{})
	dataflowTester.Subtask(tasks.ExtractDeployBuildMeta, taskData)

	dataflowTester.VerifyTable(
		models.BambooDeployBuild{},
		"./snapshot_tables/_tool_bamboo_deploy_builds.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"deploy_build_id",
			"deployment_version_name",
			"deployment_state",
			"life_cycle_state",
			"started_date",
			"queued_date",
			"executed_date",
			"finished_date",
			"reason_summary",
			"plan_key",
			"project_key",
			"can_view",
			"can_edit",
			"can_delete",
			"allowed_to_execute",
			"can_execute",
			"plan_result_key",
			"allowed_to_create_version",
			"allowed_to_set_version_status",
			"environment",
		),
	)

	// verify conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bamboo_plan_build_commits.csv", &models.BambooPlanBuildVcsRevision{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bamboo_deploy_builds.csv", &models.BambooDeployBuild{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bamboo_plans.csv", models.BambooPlan{})

	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.FlushTabler(&devops.CICDDeployment{})
	dataflowTester.Subtask(tasks.ConvertDeployBuildsToDeploymentCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdDeploymentCommit{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/cicd_deployment_commits.csv",
		IgnoreTypes:  []interface{}{common.NoPKModel{}},
		IgnoreFields: []string{"created_date", "queued_date", "started_date", "finished_date"},
	})
	dataflowTester.VerifyTableWithOptions(&devops.CicdDeploymentCommit{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/cicd_deployments.csv",
		IgnoreTypes:  []interface{}{common.NoPKModel{}},
		IgnoreFields: []string{"created_date", "queued_date", "started_date", "finished_date"},
	})
}
