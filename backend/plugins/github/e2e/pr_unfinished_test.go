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
	"sort"
	"testing"

	"github.com/apache/devlake/helpers/e2ehelper"
	"github.com/apache/devlake/plugins/github/impl"
	"github.com/apache/devlake/plugins/github/models"
	"github.com/apache/devlake/plugins/github/tasks"
	"github.com/stretchr/testify/assert"
)

// TestGithubUnfinishedPullRequests pins the set of pull requests the "Collect Pull Requests"
// detail collector re-collects on an incremental run. The seeded table covers every
// combination that decides it, and the repo/connection of another scope.
//
// Pull request 6 is the regression guard: it is closed, was never merged and its details
// were already collected, so there is nothing left to collect for it. The previous
// condition (`state != 'closed' OR merged = false`) kept it in the set on every run,
// forever, which cost one request per closed and unmerged pull request per run (#9240).
func TestGithubUnfinishedPullRequests(t *testing.T) {
	var githubPlugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", githubPlugin)

	dataflowTester.ImportCsvIntoTabler(
		"./raw_tables/_tool_github_pull_requests_unfinished.csv",
		&models.GithubPullRequest{},
	)

	var pullRequests []models.GithubPullRequest
	assert.Nil(t, dataflowTester.Dal.All(
		&pullRequests,
		tasks.UnfinishedPullRequestsClauses(134018330, 1)...,
	))

	numbers := make([]int, len(pullRequests))
	for i, pullRequest := range pullRequests {
		numbers[i] = pullRequest.Number
	}
	sort.Ints(numbers)

	assert.Equal(t, []int{
		1, // open, details never collected
		2, // open, details already collected
		4, // closed and merged on GitHub, `merged` not collected yet
		5, // closed, never merged, details never collected
	}, numbers)
	assert.NotContains(t, numbers, 6, "closed and unmerged pull requests with collected details must not be re-collected")
	assert.NotContains(t, numbers, 3, "closed and merged pull requests must not be re-collected")
	assert.NotContains(t, numbers, 7, "pull requests of another repo must not be collected")
	assert.NotContains(t, numbers, 8, "pull requests of another connection must not be collected")
}
