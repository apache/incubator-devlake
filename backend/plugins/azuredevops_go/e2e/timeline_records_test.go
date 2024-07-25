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

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/impl"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/tasks"
)

func TestAzuredevopsTimelineDataFlow(t *testing.T) {

	var azuredevops impl.Azuredevops
	dataflowTester := e2ehelper.NewDataFlowTester(t, "azuredevops_go", azuredevops)

	regexEnricher := api.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "(?i)deploy")
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "(?i)production")

	taskData := &tasks.AzuredevopsTaskData{
		Options: &tasks.AzuredevopsOptions{
			ConnectionId:   1,
			ProjectId:      "test-project",
			OrganizationId: "johndoe",
			RepositoryId:   "0d50ba13-f9ad-49b0-9b21-d29eda50ca33",
			ScopeConfig:    new(models.AzuredevopsScopeConfig),
		},
		RegexEnricher: regexEnricher,
	}

	dataflowTester.FlushTabler(&models.AzuredevopsTimelineRecord{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azuredevops_go_api_timeline_records.csv", "_raw_azuredevops_go_api_timeline_records")
	dataflowTester.Subtask(tasks.ExtractApiBuildRecordsMeta, taskData)

	dataflowTester.FlushTabler(&models.AzuredevopsBuild{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azuredevops_go_api_builds.csv", "_raw_azuredevops_go_api_builds")
	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)

	// Omit the datetime columns to avoid test failures caused by the varying precision
	// in MySQL’s datetime(3) and PostgreSQL’s 'timestamp with time zone'.
	dataflowTester.VerifyTable(
		models.AzuredevopsTimelineRecord{},
		"./snapshot_tables/_tool_azuredevops_go_timeline_records.csv",
		[]string{
			"connection_id",
			"record_id",
			"build_id",
			"parent_id",
			"type",
			"name",
			"state",
			"result",
			"change_id",
		},
	)

	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.Subtask(tasks.ConvertApiTimelineRecordsMeta, taskData)

	// Omit the datetime columns to avoid test failures caused by the varying precision
	// in MySQL’s datetime(3) and PostgreSQL’s 'timestamp with time zone'.
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks.csv",
		[]string{
			"id",
			"name",
			"pipeline_id",
			"result",
			"status",
			"original_status",
			"original_result",
			"type",
			"environment",
			"duration_sec",
			"queued_duration_sec",
			"cicd_scope_id",
		},
	)
}
