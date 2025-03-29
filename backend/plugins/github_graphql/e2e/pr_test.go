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
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	githubGraphQLTasks "github.com/apache/incubator-devlake/plugins/github_graphql/tasks"
)

func TestGithubPRsDataFlow(t *testing.T) {
	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "apache/incubator-devlake",
			GithubId:     384111310,
		},
	}

	// import raw data table
	dataflowTester.FlushTabler(&models.GithubRepo{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_graphql_prs.csv", "_raw_github_graphql_prs")
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_github_repos2.csv", &models.GithubRepo{})

	dataflowTester.Subtask(githubGraphQLTasks.ExtractPrsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubPullRequest{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/_tool_github_pull_requests.csv",
		IgnoreTypes:  []interface{}{common.NoPKModel{}},
		IgnoreFields: []string{"body"},
	})

}
