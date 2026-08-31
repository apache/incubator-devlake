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

package tasks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertGithubPullRequestMerged(t *testing.T) {
	mergedAt := time.Date(2026, 8, 14, 10, 0, 0, 0, time.UTC)

	merged, err := convertGithubPullRequest(&GraphqlQueryPr{
		DatabaseId: 1,
		Number:     1,
		State:      `MERGED`,
		MergedAt:   &mergedAt,
		ClosedAt:   &mergedAt,
	}, 1, 1)
	assert.Nil(t, err)
	assert.True(t, merged.Merged)

	closed, err := convertGithubPullRequest(&GraphqlQueryPr{
		DatabaseId: 2,
		Number:     2,
		State:      `CLOSED`,
		ClosedAt:   &mergedAt,
	}, 1, 1)
	assert.Nil(t, err)
	assert.False(t, closed.Merged)

	open, err := convertGithubPullRequest(&GraphqlQueryPr{
		DatabaseId: 3,
		Number:     3,
		State:      `OPEN`,
	}, 1, 1)
	assert.Nil(t, err)
	assert.False(t, open.Merged)
}
