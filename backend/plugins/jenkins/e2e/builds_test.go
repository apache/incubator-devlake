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
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

func TestJenkinsBuildsDataFlow(t *testing.T) {

	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	regexEnricher := api.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, `test-sub-sub-dir\/devlake.*`)
	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake`,
			JobFullName:  `Test-jenkins-dir/test-jenkins-sub-dir/test-sub-sub-dir/devlake`,
			JobPath:      `job/Test-jenkins-dir/job/test-jenkins-sub-dir/job/test-sub-sub-dir/`,
		},
		RegexEnricher: regexEnricher,
	}

	dataflowTester.FlushTabler(&models.JenkinsBuild{})
	dataflowTester.FlushTabler(&models.JenkinsBuildCommit{})
	dataflowTester.FlushTabler(&models.JenkinsStage{})

	// import raw data table
	// SELECT * FROM _raw_jenkins_api_builds INTO OUTFILE "/tmp/_raw_jenkins_api_builds.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_builds.csv", "_raw_jenkins_api_builds")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_jenkins_stages.csv", &models.JenkinsStage{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_jenkins_jobs.csv", &models.JenkinsJob{})

	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		"./snapshot_tables/_tool_jenkins_builds.csv",
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
		"./snapshot_tables/_tool_jenkins_build_commits.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"build_name",
			"commit_sha",
			"branch",
			"repo_url",
		),
	)

	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})
	dataflowTester.Subtask(tasks.EnrichApiBuildWithStagesMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		"./snapshot_tables/_tool_jenkins_builds_after_enrich.csv",
		[]string{
			"connection_id",
			"job_name",
			"duration",
			"full_name",
			"has_stages",
		},
	)

	dataflowTester.Subtask(tasks.ConvertBuildReposMeta, taskData)
	dataflowTester.Subtask(tasks.ConvertBuildsToCicdTasksMeta, taskData)

	// verify env when prod env is omitted
	dataflowTester.Subtask(tasks.ConvertBuildsToCicdTasksMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks_no_prod_regex.csv",
		e2ehelper.ColumnWithRawData(
			"environment",
		),
	)

	// continue
	_ = regexEnricher.TryAdd(devops.PRODUCTION, `test-sub-sub-dir\/devlake.*`)
	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.Subtask(tasks.ConvertBuildsToCicdTasksMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks.csv",
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
		"./snapshot_tables/cicd_pipelines.csv",
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

	dataflowTester.VerifyTable(
		devops.CiCDPipelineCommit{},
		"./snapshot_tables/cicd_pipeline_commits.csv",
		e2ehelper.ColumnWithRawData(
			"pipeline_id",
			"repo_id",
			"repo_url",
			"branch",
			"commit_sha",
			"display_title",
			"url",
		),
	)
}
