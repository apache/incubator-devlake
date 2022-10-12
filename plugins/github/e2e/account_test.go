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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
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
	dataflowTester.VerifyTable(
		models.GithubAccount{},
		"./snapshot_tables/_tool_github_account.csv",
		[]string{
			"connection_id",
			"id",
			"login",
			"name",
			"company",
			"email",
			"avatar_url",
			"url",
			"html_url",
			"type",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_account_orgs.csv", "_raw_github_api_account_orgs")
	// verify extraction
	dataflowTester.FlushTabler(&models.GithubAccountOrg{})
	dataflowTester.Subtask(tasks.ExtractAccountOrgMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubAccountOrg{},
		"./snapshot_tables/_tool_github_account_orgs.csv",
		[]string{
			"connection_id",
			"account_id",
			"org_id",
			"org_login",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify converter
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.Account{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/account.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}.CreatedAt, common.NoPKModel{}.UpdatedAt},
	})
	dataflowTester.VerifyTable(
		crossdomain.Account{},
		"./snapshot_tables/account.csv",
		[]string{
			"id",
			"email",
			"full_name",
			"user_name",
			"avatar_url",
			"organization",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
