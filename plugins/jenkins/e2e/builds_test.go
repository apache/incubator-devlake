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
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

func TestJenkinsBuildsDataFlow(t *testing.T) {

	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake`,
		},
	}

	// import raw data table
	// SELECT * FROM _raw_jenkins_api_builds INTO OUTFILE "/tmp/_raw_jenkins_api_builds.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_builds.csv", "_raw_jenkins_api_builds")

	// verify extraction
	dataflowTester.FlushTabler(&models.JenkinsBuild{})
	dataflowTester.FlushTabler(&models.JenkinsBuildCommit{})
	dataflowTester.FlushTabler(&models.JenkinsStage{})

	dataflowTester.Subtask(tasks.ExtractApiBuildsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JenkinsBuild{},
		"./snapshot_tables/_tool_jenkins_builds.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"job_name",
			"duration",
			"full_display_name",
			"estimated_duration",
			"number",
			"result",
			"timestamp",
			"start_time",
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
	dataflowTester.FlushTabler(&devops.CICDPipelineRelationship{})
	dataflowTester.Subtask(tasks.EnrichApiBuildWithStagesMeta, taskData)
	dataflowTester.Subtask(tasks.ConvertBuildsToCICDMeta, taskData)
	dataflowTester.Subtask(tasks.ConvertBuildReposMeta, taskData)

	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks.csv",
		e2ehelper.ColumnWithRawData(
			"name",
			"pipeline_id",
			"result",
			"status",
			"type",
			"environment",
			"duration_sec",
			"started_date",
			"finished_date",
		),
	)

	dataflowTester.VerifyTable(
		devops.CICDPipeline{},
		"./snapshot_tables/cicd_pipelines.csv",
		e2ehelper.ColumnWithRawData(
			"name",
			"result",
			"status",
			"type",
			"duration_sec",
			"environment",
			"created_date",
			"finished_date",
		),
	)

	dataflowTester.VerifyTable(
		devops.CICDPipelineRelationship{},
		"./snapshot_tables/cicd_pipeline_relationships.csv",
		e2ehelper.ColumnWithRawData(
			"parent_pipeline_id",
			"child_pipeline_id",
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
		),
	)
}
