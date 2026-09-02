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

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/helpers/e2ehelper"
	"github.com/apache/devlake/plugins/jira/impl"
	"github.com/apache/devlake/plugins/jira/models"
	"github.com/apache/devlake/plugins/jira/tasks"
	"github.com/stretchr/testify/assert"
)

// Board scoping in the changelog convertor is deliberate — the board is the unit of work — but it
// must not be silent, which is the complaint in issue #8834: an operator sees fewer rows in
// issue_changelogs than in the tool layer with nothing in the log to explain it.
//
// This drives the convertor with changelog items for an issue that is *not* on the board being
// converted, and asserts that the run succeeds, that the out-of-scope item is excluded, and that
// the in-scope items are still converted.
func TestIssueChangelogBoardScopeDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_changelogs.csv", &models.JiraIssueChangelogs{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_changelog_items.csv", &models.JiraIssueChangelogItems{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_statuses_for_changelog.csv", &models.JiraStatus{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_changelog.csv", &models.JiraBoardIssue{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_fields.csv", &models.JiraIssueField{})

	db := dataflowTester.Dal

	// An issue that exists with changelogs but is on no board at all — the shape the reporter
	// describes for epic-sourced and cross-referenced issues.
	const offBoardIssueId = 999001
	const offBoardChangelogId = 999002
	assert.NoError(t, db.CreateOrUpdate(&models.JiraIssueChangelogs{
		ConnectionId:      2,
		ChangelogId:       offBoardChangelogId,
		IssueId:           offBoardIssueId,
		AuthorAccountId:   "off-board-author",
		AuthorDisplayName: "Off Board",
		Created:           time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
	}))
	assert.NoError(t, db.CreateOrUpdate(&models.JiraIssueChangelogItems{
		ConnectionId: 2,
		ChangelogId:  offBoardChangelogId,
		Field:        "status",
		FromString:   "New",
		ToString:     "Closed",
	}))

	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertIssueChangelogsMeta, taskData)

	// The out-of-scope item must not appear in the domain layer...
	offBoard, err := db.Count(
		dal.From(&ticket.IssueChangelogs{}),
		dal.Where("issue_id = ?", "jira:JiraIssue:2:999001"),
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), offBoard,
		"a changelog for an issue outside the board must not be converted")

	// ...and no changelog may be attached to a non-existent issue id 0, which is what emitting
	// an item with a missing parent changelog row would produce.
	orphans, err := db.Count(
		dal.From(&ticket.IssueChangelogs{}),
		dal.Where("issue_id = ?", "jira:JiraIssue:2:0"),
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), orphans, "no changelog may point at issue id 0")

	// The in-scope changelogs are still converted, so the exclusion above is scoping and not a
	// regression that dropped everything.
	inScope, err := db.Count(dal.From(&ticket.IssueChangelogs{}))
	assert.NoError(t, err)
	assert.Greater(t, inScope, int64(0), "in-scope changelogs must still be converted")
}
