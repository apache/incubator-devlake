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
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestPrEnrichIssueDataFlow(t *testing.T) {
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
			TransformationRules: models.TransformationRules{
				PrType:               "type/(.*)$",
				PrComponent:          "component/(.*)$",
				PrBodyClosePattern:   "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)",
				IssueSeverity:        "severity/(.*)$",
				IssuePriority:        "^(highest|high|medium|low)$",
				IssueComponent:       "component/(.*)$",
				IssueTypeBug:         "^(bug|failure|error)$",
				IssueTypeIncident:    "",
				IssueTypeRequirement: "^(feat|feature|proposal|requirement)$",
			},
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_github_issues.csv", &models.GithubIssue{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_github_pull_requests.csv", models.GithubPullRequest{})

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubPullRequestIssue{})
	dataflowTester.Subtask(tasks.EnrichPullRequestIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubPullRequestIssue{},
		"./snapshot_tables/_tool_github_pull_request_issues.csv",
		[]string{
			"connection_id",
			"pull_request_id",
			"issue_id",
			"pull_request_number",
			"issue_number",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify extraction
	dataflowTester.FlushTabler(&crossdomain.PullRequestIssue{})
	dataflowTester.Subtask(tasks.ConvertPullRequestIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.PullRequestIssue{},
		"./snapshot_tables/pull_request_issues.csv",
		[]string{
			"pull_request_id",
			"issue_id",
			"pull_request_number",
			"issue_number",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
