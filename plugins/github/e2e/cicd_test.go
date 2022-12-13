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
	"github.com/apache/incubator-devlake/models/common"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestGithubCICDDataFlow(t *testing.T) {
	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			GithubId:     134018330,
			GithubTransformationRule: &models.GithubTransformationRule{
				DeploymentPattern: `deploy.*`,
				ProductionPattern: `deploywindows.*`,
			},
		},
	}

	// import raw data table
	// SELECT * FROM _raw_github_api_runs INTO OUTFILE "/tmp/_raw_github_api_runs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_runs.csv", "_raw_github_api_runs")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubRun{})
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})

	dataflowTester.Subtask(tasks.ExtractRunsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubRun{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_runs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.Subtask(tasks.ConvertRunsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDPipeline{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipelines.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&devops.CiCDPipelineCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipeline_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import raw data table
	// SELECT * FROM _raw_github_api_jobs INTO OUTFILE "/tmp/_raw_github_api_jobs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_jobs.csv", "_raw_github_api_jobs")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubJob{})
	dataflowTester.FlushTabler(&devops.CICDTask{})

	dataflowTester.Subtask(tasks.ExtractJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubJob{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_jobs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDTask{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_tasks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
