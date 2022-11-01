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
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestAccountDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)

	githubRepository := &models.GithubRepo{
		GithubId: 134018330,
	}
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_accounts.csv", "_raw_github_api_accounts")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubAccount{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_account.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_account_orgs.csv", "_raw_github_api_account_orgs")
	// verify extraction
	dataflowTester.FlushTabler(&models.GithubAccountOrg{})
	dataflowTester.Subtask(tasks.ExtractAccountOrgMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubAccountOrg{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_account_orgs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// ConvertAccountsMeta only convert the account in this repo
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_github_repo_accounts.csv", &models.GithubRepoAccount{})

	// verify converter
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.Account{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/account.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
