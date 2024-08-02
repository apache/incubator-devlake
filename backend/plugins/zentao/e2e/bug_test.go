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
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

func TestZentaoBugDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    1,
			ScopeConfig: &models.ZentaoScopeConfig{
				TypeMappings: map[string]string{
					"codeerror": "CODE_ERROR",
				},
				BugStatusMappings: map[string]string{
					"active": ticket.DONE,
				},
			},
		},
		Bugs:         map[int64]struct{}{},
		AccountCache: tasks.NewAccountCache(dataflowTester.Dal, 1),
		ApiClient:    getFakeAPIClient(),
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_bugs.csv",
		"_raw_zentao_api_bugs")

	// verify extraction
	dataflowTester.FlushTabler(&models.ZentaoBug{})
	dataflowTester.Subtask(tasks.ExtractBugMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoBug{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_bugs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify conversion
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.FlushTabler(&ticket.SprintIssue{})
	dataflowTester.FlushTabler(&ticket.IssueAssignee{})
	dataflowTester.Subtask(tasks.ConvertBugMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.Issue{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/issues_bug.csv",
		IgnoreTypes:  []interface{}{common.NoPKModel{}},
		IgnoreFields: []string{"original_project"},
	})
	dataflowTester.VerifyTableWithOptions(&ticket.BoardIssue{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/board_issues_bug.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.VerifyTableWithOptions(ticket.IssueAssignee{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/bug_issue_assignees.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
