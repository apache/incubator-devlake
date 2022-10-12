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

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestMilestoneDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)
	githubRepository := &models.GithubRepo{
		ConnectionId: 1,
		GithubId:     134018330,
		CreatedDate: func() time.Time {
			createdTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
			return createdTime
		}(),
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

	dataflowTester.FlushTabler(&models.GithubMilestone{})
	dataflowTester.FlushTabler(&models.GithubIssue{})

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_milestones.csv", "_raw_"+tasks.RAW_MILESTONE_TABLE)
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_issues.csv", "_raw_"+tasks.RAW_ISSUE_TABLE)

	dataflowTester.Subtask(tasks.ExtractApiIssuesMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractMilestonesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubMilestone{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/_tool_github_milestones.csv",
	})

	dataflowTester.FlushTabler(&ticket.Sprint{})
	dataflowTester.FlushTabler(&ticket.BoardSprint{})
	dataflowTester.FlushTabler(&ticket.SprintIssue{})

	dataflowTester.Subtask(tasks.ConvertMilestonesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.Sprint{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/sprints.csv",
	})
	dataflowTester.VerifyTableWithOptions(&ticket.BoardSprint{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/board_sprint.csv",
	})
	dataflowTester.VerifyTableWithOptions(&ticket.SprintIssue{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/sprint_issue.csv",
	})
}
