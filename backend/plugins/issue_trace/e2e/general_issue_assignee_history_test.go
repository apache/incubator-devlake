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
	"github.com/apache/incubator-devlake/plugins/issue_trace/impl"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
	"github.com/apache/incubator-devlake/plugins/issue_trace/tasks"
)

func TestConvertIssueAssigneeHistory(t *testing.T) {
	var plugin impl.IssueTrace
	dataflowTester := e2ehelper.NewDataFlowTester(t, "issue_trace", plugin)

	dataflowTester.ImportCsvIntoTabler("./raw_tables/board_issues.csv", &ticket.BoardIssue{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issues.csv", &ticket.Issue{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/issue_changelogs.csv", &ticket.IssueChangelogs{})

	dataflowTester.FlushTabler(models.IssueAssigneeHistory{})

	dataflowTester.Subtask(tasks.ConvertIssueAssigneeHistoryMeta, TaskData)

	dataflowTester.VerifyTable(
		models.IssueAssigneeHistory{},
		"./snapshot_tables/issue_assignee_history.csv",
		[]string{
			"issue_id",
			"assignee",
			"start_date",
		},
	)

}
