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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestGithubCICDRunDataFlow(t *testing.T) {
	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)
	regexEnricher := helper.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "CodeQL.*")
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
		},
		RegexEnricher: regexEnricher,
	}

	// import raw data table
	// SELECT * FROM _raw_github_api_runs INTO OUTFILE "/tmp/_raw_github_api_runs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_runs.csv", "_raw_github_api_runs")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_github_repos.csv", &models.GithubRepo{})

	// verify when production regex is omitted
	dataflowTester.FlushTabler(&models.GithubRun{})
	dataflowTester.Subtask(tasks.ExtractRunsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubRun{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_runs_no_prod_regex.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubRun{})
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "CodeQL.*")
	dataflowTester.Subtask(tasks.ExtractRunsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubRun{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_runs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})
	dataflowTester.Subtask(tasks.ConvertRunsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDPipeline{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipelines.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(&devops.CiCDPipelineCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipeline_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
