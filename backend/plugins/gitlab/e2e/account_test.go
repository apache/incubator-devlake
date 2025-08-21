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
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabAccountDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)
	apiClient := &api.ApiClient{}
	apiClient.Setup(
		"https://gitlab.com",
		make(map[string]string),
		time.Hour,
	)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12345678,
			ScopeConfig:  new(models.GitlabScopeConfig),
		},

		ApiClient: &api.ApiAsyncClient{
			ApiClient: apiClient,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_users.csv",
		"_raw_gitlab_api_users")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabAccount{},
		"./snapshot_tables/_tool_gitlab_accounts.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"gitlab_id",
			"username",
			"email",
			"name",
			"state",
			"membership_state",
			"avatar_url",
			"web_url",
		),
	)

	// verify conversion
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.Account{},
		"./snapshot_tables/accounts.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"email",
			"full_name",
			"user_name",
			"avatar_url",
			"organization",
			"created_date",
			"status",
		),
	)
}

func TestGitlabAccountDataFlowUsersApi(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)
	apiClient := &api.ApiClient{}
	apiClient.Setup(
		"https://custom.gitlab.com",
		make(map[string]string),
		time.Hour,
	)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12345678,
			ScopeConfig:  new(models.GitlabScopeConfig),
		},

		ApiClient: &api.ApiAsyncClient{
			ApiClient: apiClient,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_users_direct_api.csv",
		"_raw_gitlab_api_users")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabAccount{},
		"./snapshot_tables/_tool_gitlab_accounts_direct_api.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"gitlab_id",
			"username",
			"email",
			"name",
			"state",
			"avatar_url",
			"web_url",
			"created_user_at",
		),
	)

	// verify conversion
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.Account{},
		"./snapshot_tables/accounts_direct_api.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"email",
			"full_name",
			"user_name",
			"avatar_url",
			"organization",
			"created_date",
			"status",
		),
	)
}
