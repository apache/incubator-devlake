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

package parser

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildTestRepo creates a repository on disk with a known shape: 3 commits on
// the default branch, one lightweight tag and a second branch which is checked
// out at the end, so HEAD does not point at the default branch.
func buildTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	repo, err := gogit.PlainInit(dir, false)
	require.NoError(t, err)
	wt, err := repo.Worktree()
	require.NoError(t, err)

	sig := func(n int) *object.Signature {
		return &object.Signature{
			Name:  "DevLake Test",
			Email: "test@devlake.apache.org",
			When:  time.Date(2026, 1, 1, 0, 0, n, 0, time.UTC),
		}
	}

	var last plumbing.Hash
	for i := 1; i <= 3; i++ {
		name := filepath.Join(dir, "file.txt")
		require.NoError(t, os.WriteFile(name, []byte{byte('a' + i)}, 0600))
		_, err = wt.Add("file.txt")
		require.NoError(t, err)
		last, err = wt.Commit("commit", &gogit.CommitOptions{Author: sig(i)})
		require.NoError(t, err)
	}

	_, err = repo.CreateTag("v1.0.0", last, nil)
	require.NoError(t, err)

	require.NoError(t, wt.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("feature"),
		Create: true,
	}))

	return dir
}

func TestGogitRepoCollectorCounts(t *testing.T) {
	dir := buildTestRepo(t)

	// store and logger are only used by the Collect* subtasks, not by the
	// counting helpers exercised here.
	collector, err := NewGogitRepoCollector(dir, "test-repo", nil, nil)
	require.Nil(t, err)

	ctx := context.Background()

	commits, cErr := collector.CountCommits(ctx)
	assert.NoError(t, cErr)
	assert.Equal(t, 3, commits)

	tags, tErr := collector.CountTags(ctx)
	assert.NoError(t, tErr)
	assert.Equal(t, 1, tags)

	// HEAD is on "feature", so only the default branch is counted.
	branches, bErr := collector.CountBranches(ctx)
	assert.NoError(t, bErr)
	assert.Equal(t, 1, branches)
}

func TestGogitRepoCollectorCountsHonourCancelledContext(t *testing.T) {
	dir := buildTestRepo(t)

	collector, err := NewGogitRepoCollector(dir, "test-repo", nil, nil)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, cErr := collector.CountCommits(ctx)
	assert.ErrorIs(t, cErr, context.Canceled)

	_, tErr := collector.CountTags(ctx)
	assert.ErrorIs(t, tErr, context.Canceled)

	_, bErr := collector.CountBranches(ctx)
	assert.ErrorIs(t, bErr, context.Canceled)
}

func TestNewGogitRepoCollectorRejectsNonRepo(t *testing.T) {
	_, err := NewGogitRepoCollector(t.TempDir(), "test-repo", nil, nil)
	assert.NotNil(t, err)
}

