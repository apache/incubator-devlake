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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/impl"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/tasks"
)

func TestAzuredevopsRrCommitDataFlow(t *testing.T) {

	var azuredevops impl.Azuredevops
	dataflowTester := e2ehelper.NewDataFlowTester(t, "azuredevops_go", azuredevops)

	taskData := &tasks.AzuredevopsTaskData{
		Options: &tasks.AzuredevopsOptions{
			ConnectionId:   1,
			ProjectId:      "test-project",
			OrganizationId: "johndoe",
			RepositoryId:   "0d50ba13-f9ad-49b0-9b21-d29eda50ca33",
			ScopeConfig:    new(models.AzuredevopsScopeConfig),
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azuredevops_go_api_pull_requests.csv",
		"_raw_azuredevops_go_api_pull_requests")
	// verify extraction
	dataflowTester.FlushTabler(&models.AzuredevopsPullRequest{})
	dataflowTester.FlushTabler(&models.AzuredevopsPrLabel{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestsMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&models.AzuredevopsPullRequest{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_azuredevops_go_pull_requests_for_pr_commit_test.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_azuredevops_go_api_pull_request_commits.csv",
		"_raw_azuredevops_go_api_pull_request_commits")

	// verify extraction
	dataflowTester.FlushTabler(&models.AzuredevopsCommit{})
	dataflowTester.FlushTabler(&models.AzuredevopsRepoCommit{})
	dataflowTester.FlushTabler(&models.AzuredevopsPrCommit{})
	dataflowTester.Subtask(tasks.ExtractApiPullRequestCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(models.AzuredevopsCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_azuredevops_go_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.VerifyTableWithOptions(models.AzuredevopsRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_azuredevops_go_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.VerifyTableWithOptions(models.AzuredevopsPrCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_azuredevops_go_pull_request_commits.csv",
		IgnoreTypes: []interface{}{common.Model{}},
	})

	// verify conversion
	dataflowTester.FlushTabler(&code.PullRequestCommit{})
	dataflowTester.Subtask(tasks.ConvertApiPrCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(code.PullRequestCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/pull_request_commits.csv",
		IgnoreTypes: []interface{}{common.Model{}},
	})

	// verify conversion
	dataflowTester.FlushTabler(&code.Commit{})
	dataflowTester.FlushTabler(&code.RepoCommit{})
	dataflowTester.Subtask(tasks.ConvertCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(code.Commit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/commits.csv",
		IgnoreTypes: []interface{}{common.Model{}},
	})

	dataflowTester.VerifyTable(
		code.RepoCommit{},
		"./snapshot_tables/repo_commits.csv",
		e2ehelper.ColumnWithRawData(
			"repo_id",
			"commit_sha",
		),
	)
}
