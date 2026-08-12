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

package api

import (
	"net/url"
	"testing"

	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/stretchr/testify/assert"
)

func TestGroupScopeDuplicateRows_Empty(t *testing.T) {
	assert.Empty(t, groupScopeDuplicateRows(nil))
	assert.Empty(t, groupScopeDuplicateRows([]scopeDuplicateRow{}))
}

func TestGroupScopeDuplicateRows_GroupsConnections(t *testing.T) {
	rows := []scopeDuplicateRow{
		{BitbucketId: "o/a", HTMLUrl: "https://bitbucket.org/o/a", ConnectionId: 1, ConnectionName: "Bitbucket Production"},
		{BitbucketId: "o/a", HTMLUrl: "https://bitbucket.org/o/a", ConnectionId: 2, ConnectionName: "Bitbucket Staging"},
		{BitbucketId: "o/b", HTMLUrl: "https://bitbucket.org/o/b", ConnectionId: 3, ConnectionName: "Other"},
	}

	got := groupScopeDuplicateRows(rows)
	assert.Equal(t, []ScopeDuplicateGroup{
		{
			BitbucketId: "o/a",
			HTMLUrl:     "https://bitbucket.org/o/a",
			FullName:    "o/a",
			Connections: []ScopeDuplicateConnection{
				{ConnectionId: 1, ConnectionName: "Bitbucket Production"},
				{ConnectionId: 2, ConnectionName: "Bitbucket Staging"},
			},
		},
		{
			BitbucketId: "o/b",
			HTMLUrl:     "https://bitbucket.org/o/b",
			FullName:    "o/b",
			Connections: []ScopeDuplicateConnection{
				{ConnectionId: 3, ConnectionName: "Other"},
			},
		},
	}, got)
}

func TestGroupScopeDuplicateRows_DedupesSameConnection(t *testing.T) {
	rows := []scopeDuplicateRow{
		{BitbucketId: "o/a", ConnectionId: 1, ConnectionName: "Prod"},
		{BitbucketId: "o/a", ConnectionId: 1, ConnectionName: "Prod"},
	}

	got := groupScopeDuplicateRows(rows)
	assert.Len(t, got, 1)
	assert.Equal(t, []ScopeDuplicateConnection{
		{ConnectionId: 1, ConnectionName: "Prod"},
	}, got[0].Connections)
}

func TestGroupScopeDuplicateRows_FillsMissingHTMLUrl(t *testing.T) {
	rows := []scopeDuplicateRow{
		{BitbucketId: "o/a", ConnectionId: 1, ConnectionName: "A"},
		{BitbucketId: "o/a", HTMLUrl: "https://bitbucket.org/o/a", ConnectionId: 2, ConnectionName: "B"},
	}

	got := groupScopeDuplicateRows(rows)
	assert.Len(t, got, 1)
	assert.Equal(t, "https://bitbucket.org/o/a", got[0].HTMLUrl)
	assert.Equal(t, "o/a", got[0].FullName)
}

func TestParseScopeDuplicateQuery(t *testing.T) {
	input := &plugin.ApiResourceInput{Query: url.Values{}}
	connId, ids, err := parseScopeDuplicateQuery(input)
	assert.Nil(t, err)
	assert.Nil(t, connId)
	assert.Empty(t, ids)

	input = &plugin.ApiResourceInput{Query: url.Values{
		"connectionId": []string{"3"},
		"bitbucketIds": []string{"o/a, o/b,o/c"},
	}}
	connId, ids, err = parseScopeDuplicateQuery(input)
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), *connId)
	assert.Equal(t, []string{"o/a", "o/b", "o/c"}, ids)

	input = &plugin.ApiResourceInput{Query: url.Values{
		"bitbucketIds": []string{"o/a"},
	}}
	_, _, err = parseScopeDuplicateQuery(input)
	assert.Contains(t, err.Error(), "connectionId is required when bitbucketIds is provided")

	input = &plugin.ApiResourceInput{Query: url.Values{
		"connectionId": []string{"abc"},
	}}
	_, _, err = parseScopeDuplicateQuery(input)
	assert.Contains(t, err.Error(), "invalid connectionId")
}
