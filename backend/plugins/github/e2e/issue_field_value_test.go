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

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

// Extraction of issue field values, and the scope config mapping that carries them into the
// domain issue.
func TestIssueFieldValueDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
			ScopeConfig: &models.GithubScopeConfig{
				// Mapped by field name. "Effort" and "Target date" are two of the four fields
				// GitHub preconfigures for every organization.
				IssueFieldPriority:   "Priority",
				IssueFieldComponent:  "Squad",
				IssueFieldStoryPoint: "Effort",
				IssueFieldDueDate:    "Target date",
			},
		},
	}

	dataflowTester.ImportCsvIntoRawTable(
		"./raw_tables/_raw_github_api_issue_field_values.csv",
		"_raw_github_api_issue_field_values")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubIssueFieldValue{})
	dataflowTester.Subtask(tasks.ExtractApiIssueFieldValuesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		models.GithubIssueFieldValue{},
		e2ehelper.TableOptions{
			CSVRelPath: "./snapshot_tables/_tool_github_issue_field_values.csv",
			TargetFields: []string{
				"connection_id",
				"issue_id",
				"field_id",
				"field_name",
				"data_type",
				"value",
				"raw_value",
				"option_color",
			}},
	)

	// verify the mapping onto the domain issue
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_github_issues.csv", &models.GithubIssue{})
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.Subtask(tasks.ConvertIssuesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		&ticket.Issue{},
		e2ehelper.TableOptions{
			CSVRelPath: "./snapshot_tables/issues_with_field_values.csv",
			TargetFields: []string{
				"id",
				"issue_key",
				"priority",
				"component",
				"story_point",
				"due_date",
			},
			IgnoreTypes: []interface{}{},
		},
	)
}
