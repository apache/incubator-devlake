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
	githubGraphQLTasks "github.com/apache/incubator-devlake/plugins/github_graphql/tasks"
)

func TestGithubDeploymentDataFlow(t *testing.T) {
	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)
	regexEnricher := helper.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "github-pages")
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "facebook/OpenBIC",
			GithubId:     335709078,
		},
		RegexEnricher: regexEnricher,
	}

	// import raw data table
	dataflowTester.FlushTabler(&models.GithubDeployment{})
	dataflowTester.FlushTabler(&models.GithubRepo{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_graphql_deployment.csv", "_raw_github_graphql_deployment")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_github_repos2.csv", &models.GithubRepo{})

	dataflowTester.Subtask(githubGraphQLTasks.ExtractDeploymentsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubDeployment{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_deployments.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify convertor
	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.FlushTabler(&devops.CICDDeployment{})
	dataflowTester.Subtask(tasks.ConvertDeploymentsMeta, taskData)
	dataflowTester.VerifyTable(&devops.CicdDeploymentCommit{},
		"./snapshot_tables/cicd_deployment_commits.csv",
		[]string{
			"cicd_scope_id",
			"cicd_deployment_id",
			"name",
			"result",
			"status",
			"original_status",
			"environment",
			"original_environment",
			"created_date",
			"queued_date",
			"started_date",
			"finished_date",
			"commit_sha",
			"commit_msg",
			"ref_name",
			"repo_id",
			"repo_url",
			"prev_success_deployment_commit_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.VerifyTable(&devops.CICDDeployment{},
		"./snapshot_tables/cicd_deployments.csv",
		[]string{
			"cicd_scope_id",
			"name",
			"result",
			"status",
			"original_status",
			"environment",
			"created_date",
			"queued_date",
			"started_date",
			"finished_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
