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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/stretchr/testify/assert"
)

func TestFindMissingGithubIssues(t *testing.T) {
	requestedIssues := map[int]missingGithubIssueRef{
		17: {
			ConnectionId: 1,
			GithubId:     1700,
			Number:       17,
			RawDataOrigin: common.RawDataOrigin{
				RawDataTable: "_raw_github_graphql_issues",
				RawDataId:    10,
			},
		},
		18: {
			ConnectionId: 1,
			GithubId:     1800,
			Number:       18,
		},
	}

	resolvedIssues := []GraphqlQueryIssue{
		{DatabaseId: 1800, Number: 18},
	}

	missingIssues := findMissingGithubIssues(requestedIssues, resolvedIssues)

	if assert.Len(t, missingIssues, 1) {
		assert.Equal(t, 17, missingIssues[0].Number)
		assert.Equal(t, 1700, missingIssues[0].GithubId)
		assert.Equal(t, uint64(10), missingIssues[0].RawDataOrigin.RawDataId)
	}
}

func TestFindMissingGithubIssuesSkipsZeroValueResponses(t *testing.T) {
	requestedIssues := map[int]missingGithubIssueRef{
		17: {Number: 17},
		18: {Number: 18},
	}

	resolvedIssues := []GraphqlQueryIssue{
		{},
		{DatabaseId: 1800, Number: 18},
	}

	missingIssues := findMissingGithubIssues(requestedIssues, resolvedIssues)

	if assert.Len(t, missingIssues, 1) {
		assert.Equal(t, 17, missingIssues[0].Number)
	}
}
