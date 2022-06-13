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
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ProjectId: 3472737,
		},
	}

	// import raw data table
	dataflowTester.MigrateRawTableAndFlush("_raw_gitlab_api_project")
	dataflowTester.ImportCsv("./tables/_raw_gitlab_api_projects.csv", "_raw_gitlab_api_project")

	// verify extraction
	dataflowTester.MigrateTableAndFlush(&models.GithubIssue{})
	dataflowTester.Subtask(tasks.ExtractProjectMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		"_tool_gitlab_projects",
		"tables/_tool_gitlab_projects.csv",
		[]string{"gitlab_id"},
		[]string{
			"name",
			"description",
			"default_branch",
			"path_with_namespace",
			"web_url",
			"creator_id",
			"visibility",
			"open_issues_count",
			"star_count",
			"forked_from_project_id",
			"forked_from_project_web_url",
			"created_date",
			"updated_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTable("repos")
	dataflowTester.Subtask(tasks.ConvertProjectMeta, taskData)
	dataflowTester.CreateSnapshotOrVerify(
		"repos",
		"tables/repos.csv",
		[]string{"id"},
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"name",
			"url",
			"description",
			"owner_id",
			"language",
			"forked_from",
			"created_date",
			"updated_date",
			"deleted",
		},
	)
}
