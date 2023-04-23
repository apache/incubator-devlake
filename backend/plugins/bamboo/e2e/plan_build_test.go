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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/impl"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

func TestBambooPlanBuildDataFlow(t *testing.T) {

	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)

	taskData := &tasks.BambooTaskData{
		Options: &models.BambooOptions{
			ConnectionId: 3,
			ProjectKey:   "TEST1",
			BambooTransformationRule: &models.BambooTransformationRule{
				DeploymentPattern: "(?i)compile",
				ProductionPattern: "(?i)compile",
			},
		},
		RegexEnricher: helper.NewRegexEnricher(),
	}
	taskData.RegexEnricher.TryAdd(devops.DEPLOYMENT, taskData.Options.DeploymentPattern)
	taskData.RegexEnricher.TryAdd(devops.PRODUCTION, taskData.Options.ProductionPattern)
	// import raw data table
	// SELECT * FROM _raw_bamboo_api_job_build INTO OUTFILE "/tmp/_raw_bamboo_api_job_build.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bamboo_api_plan_build.csv", "_raw_bamboo_api_plan_build")

	// verify extraction
	dataflowTester.FlushTabler(&models.BambooPlanBuild{})
	dataflowTester.Subtask(tasks.ExtractPlanBuildMeta, taskData)
	dataflowTester.VerifyTable(
		models.BambooPlanBuild{},
		"./snapshot_tables/_tool_bamboo_plan_builds.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"plan_build_key",
			"expand",
			"number",
			"build_number",
			"plan_name",
			"plan_key",
			"project_name",
			"project_key",
			"build_result_key",
			"life_cycle_state",
			"build_started_time",
			"pretty_build_started_time",
			"build_completed_time",
			"build_completed_date",
			"pretty_build_completed_time",
			"build_duration_in_seconds",
			"build_duration",
			"build_duration_description",
			"build_relative_time",
			"vcs_revision_key",
			"build_test_summary",
			"successful_test_count",
			"failed_test_count",
			"quarantined_test_count",
			"skipped_test_count",
			"continuable",
			"once_off",
			"restartable",
			"not_run_yet",
			"build_reason",
			"reason_summary",
			"state",
			"build_state",
			"plan_result_key",
			"type",
			"environment",
		),
	)

	// verify conversion
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})

	dataflowTester.Subtask(tasks.ConvertPlanBuildsMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&devops.CICDPipeline{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipelines.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
