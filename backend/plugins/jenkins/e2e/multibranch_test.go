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
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"

	api "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func TestJenkinsMultibranchDataFlow(t *testing.T) {
	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	regexEnricher := api.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, `test-sub-sub-dir\/devlake.*`)
	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake-jenkins`,
			JobFullName:  `github_org/devlake-jenkins`,
			JobPath:      ``,
			Class:        `org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject`,
			ScopeConfig:  new(models.JenkinsScopeConfig),
		},
		Connection: &models.JenkinsConnection{
			JenkinsConn: models.JenkinsConn{
				RestConnection: api.RestConnection{
					Endpoint: "https://1457-62-195-68-26.ngrok-free.app",
				},
			},
		},
		RegexEnricher: regexEnricher,
	}

	// This replaces tasks.CollectApiJobsMeta
	dataflowTester.FlushTabler(&models.JenkinsJob{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_jobs_multibranch.csv", "_raw_jenkins_api_jobs")

	dataflowTester.Subtask(tasks.ExtractApiJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.JenkinsJob{}, e2ehelper.TableOptions{
		CSVRelPath:  "./raw_tables/_tool_jenkins_jobs_multibranch.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdScope{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_scopes_multibranch.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// This replaces tasks.CollectApiBuildsMeta
	dataflowTester.FlushTabler(&models.JenkinsBuild{})
	dataflowTester.FlushTabler(&models.JenkinsBuildCommit{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_builds_multibranch.csv", "_raw_jenkins_api_builds")

	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		"./snapshot_tables/_tool_jenkins_builds_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"job_name",
			"job_path",
			"duration",
			"full_name",
			"estimated_duration",
			"number",
			"result",
			"timestamp",
			"start_time",
			"has_stages",
		),
	)
	dataflowTester.VerifyTable(
		models.JenkinsBuildCommit{},
		"./snapshot_tables/_tool_jenkins_build_commits_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"build_name",
			"commit_sha",
			"branch",
			"repo_url",
		),
	)

	// This replaces tasks.CollectApiStagesMeta
	dataflowTester.FlushRawTable("_raw_jenkins_api_stages")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_stages_multibranch.csv", "_raw_jenkins_api_stages")

	dataflowTester.FlushTabler(&models.JenkinsStage{})
	dataflowTester.Subtask(tasks.ExtractApiStagesMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsStage{},
		"./snapshot_tables/_tool_jenkins_stages_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"build_name",
			"name",
			"exec_node",
			"status",
			"start_time_millis",
			"duration_millis",
			"pause_duration_millis",
			"type",
		),
	)

	dataflowTester.Subtask(tasks.EnrichApiBuildWithStagesMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		"./snapshot_tables/_tool_jenkins_builds_multibranch_after_enrich.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"job_name",
			"job_path",
			"duration",
			"full_name",
			"estimated_duration",
			"number",
			"result",
			"timestamp",
			"start_time",
			"has_stages",
		),
	)

	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})

	dataflowTester.Subtask(tasks.ConvertBuildsToCicdTasksMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"name",
			"pipeline_id",
			"result",
			"status",
			"original_result",
			"original_status",
			"type",
			"environment",
			"duration_sec",
			"started_date",
			"finished_date",
			"cicd_scope_id",
		),
	)
	dataflowTester.VerifyTable(
		devops.CICDPipeline{},
		"./snapshot_tables/cicd_pipelines_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"name",
			"result",
			"status",
			"original_result",
			"original_status",
			"type",
			"duration_sec",
			"environment",
			"created_date",
			"finished_date",
			"cicd_scope_id",
		),
	)
	dataflowTester.Subtask(tasks.ConvertBuildReposMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CiCDPipelineCommit{},
		"./snapshot_tables/cicd_pipeline_commits_multibranch.csv",
		e2ehelper.ColumnWithRawData(
			"pipeline_id",
			"repo_id",
			"repo_url",
			"branch",
			"commit_sha",
		),
	)

}
