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
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
)

func TestIssueRelationshipDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_issue_relationships.csv", "_raw_jira_api_issues")

	dataflowTester.FlushTabler(&models.JiraIssueRelationship{})
	dataflowTester.Subtask(tasks.ExtractIssuesMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&models.JiraIssueRelationship{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_jira_issue_relationships.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues.csv", &models.JiraBoardIssue{})

	// verify issue conversion
	dataflowTester.FlushTabler(&ticket.IssueRelationship{})
	dataflowTester.Subtask(tasks.ConvertIssueRelationshipsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.IssueRelationship{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_relationships.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
