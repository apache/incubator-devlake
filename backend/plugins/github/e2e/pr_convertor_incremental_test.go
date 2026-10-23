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
	"context"
	"sort"
	"testing"
	"time"

	"github.com/apache/devlake/core/dal"
	coreModels "github.com/apache/devlake/core/models"
	"github.com/apache/devlake/core/models/domainlayer/code"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/core/runner"
	"github.com/apache/devlake/core/utils"
	"github.com/apache/devlake/helpers/e2ehelper"
	contextimpl "github.com/apache/devlake/impls/context"
	"github.com/apache/devlake/plugins/github/impl"
	"github.com/apache/devlake/plugins/github/models"
	"github.com/apache/devlake/plugins/github/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runSubtaskIncrementally executes one subtask with the subtask state left in place and a
// sync policy that permits incremental mode, which DataFlowTester.Subtask cannot do: it
// flushes _devlake_subtask_state and forces FullSync on every call
// (helpers/e2ehelper/data_flow_tester.go). Same approach as
// plugins/youtrack/e2e/stateful_scenario_test.go, with the subtask and plugin names passed
// in their own positions so the state key matches the one the subtask itself computes.
func runSubtaskIncrementally(
	t *testing.T,
	dataflowTester *e2ehelper.DataFlowTester,
	subtaskMeta plugin.SubTaskMeta,
	taskData interface{},
) {
	t.Helper()
	subtaskCtx := contextimpl.NewStandaloneSubTaskContext(
		context.Background(),
		runner.CreateBasicRes(dataflowTester.Cfg, dataflowTester.Log, dataflowTester.Db),
		subtaskMeta.Name,
		taskData,
		dataflowTester.Name,
		&coreModels.SyncPolicy{},
	)
	require.NoError(t, subtaskMeta.EntryPoint(subtaskCtx))
}

// TestGithubPrConvertorIncremental covers the incremental filter of "Convert Pull Requests".
//
// `stateManager.GetSince()` is the start time of the previous successful run of the subtask
// (`_devlake_subtask_state.prev_started_at`), so the filter has to select the tool records
// DevLake wrote since then, i.e. the ones whose own `updated_at` is at or after it.
//
// Pull request 29 is the regression guard: GitHub last touched it in 2019, DevLake rewrote
// its row after the previous run (`state` and `merged` were corrected), and the domain record
// still has to be updated. Filtering on `github_updated_at` compared that 2019 timestamp
// against a 2026 run time and skipped the record on this run and on every run after it, which
// left the domain record missing or stuck on its old status (#9150).
func TestGithubPrConvertorIncremental(t *testing.T) {
	var githubPlugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", githubPlugin)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
		},
	}

	prevStartedAt := time.Date(2026, 9, 16, 16, 24, 0, 0, time.UTC)
	// GitHub's timestamp for pull request 29, older than every DevLake run.
	collectedLongAgo := time.Date(2019, 8, 8, 8, 31, 19, 0, time.UTC)
	writtenAfterPrevRun := time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC)
	writtenBeforePrevRun := time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC)

	dataflowTester.FlushTabler(&models.GithubPullRequest{})
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.FlushTabler(&coreModels.SubtaskState{})

	// A pull request GitHub last updated in 2019 whose row DevLake rewrote after the previous
	// run, one updated on GitHub after the previous run, and one untouched since before it.
	seedPullRequest(t, dataflowTester, 101, 29, "closed", true, collectedLongAgo, writtenAfterPrevRun)
	seedPullRequest(t, dataflowTester, 102, 30, "closed", true, writtenAfterPrevRun, writtenAfterPrevRun)
	seedPullRequest(t, dataflowTester, 103, 31, "open", false, writtenBeforePrevRun, writtenBeforePrevRun)

	// Previous successful run of this subtask. PrevConfig is left empty so the subtask is not
	// treated as reconfigured, which would force full sync.
	require.NoError(t, dataflowTester.Dal.CreateOrUpdate(&coreModels.SubtaskState{
		Plugin:        "github",
		Subtask:       tasks.ConvertPullRequestsMeta.Name,
		Params:        utils.ToJsonString(tasks.GithubApiParams{ConnectionId: 1, Name: "panjf2000/ants"}),
		PrevStartedAt: &prevStartedAt,
	}))

	runSubtaskIncrementally(t, dataflowTester, tasks.ConvertPullRequestsMeta, taskData)

	var domainPullRequests []code.PullRequest
	require.NoError(t, dataflowTester.Dal.All(&domainPullRequests))
	converted := make([]int, len(domainPullRequests))
	statuses := map[int]string{}
	for i, domainPullRequest := range domainPullRequests {
		converted[i] = domainPullRequest.PullRequestKey
		statuses[domainPullRequest.PullRequestKey] = domainPullRequest.Status
	}
	sort.Ints(converted)

	assert.Equal(t, []int{
		29, // rewritten by DevLake after the previous run, last updated on GitHub in 2019
		30, // updated on GitHub after the previous run
	}, converted)
	assert.Equal(t, code.MERGED, statuses[29], "the corrected status has to reach the domain layer")
	assert.NotContains(t, converted, 31, "records untouched since before the previous run must be skipped")
}

// seedPullRequest inserts a tool-layer pull request and then forces `updated_at`, which gorm
// would otherwise set to the time of the insert.
func seedPullRequest(
	t *testing.T,
	dataflowTester *e2ehelper.DataFlowTester,
	githubId int,
	number int,
	state string,
	merged bool,
	githubUpdatedAt time.Time,
	updatedAt time.Time,
) {
	t.Helper()
	pullRequest := &models.GithubPullRequest{
		ConnectionId:    1,
		GithubId:        githubId,
		RepoId:          134018330,
		Number:          number,
		State:           state,
		Merged:          merged,
		GithubCreatedAt: githubUpdatedAt,
		GithubUpdatedAt: githubUpdatedAt,
	}
	if merged {
		pullRequest.MergedAt = &githubUpdatedAt
		pullRequest.ClosedAt = &githubUpdatedAt
	}
	require.NoError(t, dataflowTester.Dal.CreateOrUpdate(pullRequest))
	require.NoError(t, dataflowTester.Dal.UpdateColumn(
		&models.GithubPullRequest{},
		"updated_at",
		updatedAt,
		dal.Where("connection_id = ? AND github_id = ?", 1, githubId),
	))
}
