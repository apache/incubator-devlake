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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildIncidentsQuery_FirstPage(t *testing.T) {
	q := buildIncidentsQuery(250, "")
	assert.Equal(t, "250", q.Get("page_size"))
	assert.False(t, q.Has("after"), "first page must not send an after cursor")
}

func TestBuildIncidentsQuery_SubsequentPage(t *testing.T) {
	q := buildIncidentsQuery(250, "01FCNDV6P870EA6S7TK1DSYDG0")
	assert.Equal(t, "250", q.Get("page_size"))
	assert.Equal(t, "01FCNDV6P870EA6S7TK1DSYDG0", q.Get("after"))
}

func TestCollectedIncidentsEnvelope(t *testing.T) {
	body := []byte(`{
		"incidents": [{"id": "inc_1"}, {"id": "inc_2"}],
		"pagination_meta": {"after": "cursor-abc", "page_size": 2}
	}`)
	rawResult := collectedIncidents{}
	require.NoError(t, json.Unmarshal(body, &rawResult))
	require.Len(t, rawResult.Incidents, 2)
	require.NotNil(t, rawResult.PaginationMeta)
	require.NotNil(t, rawResult.PaginationMeta.After)
	assert.Equal(t, "cursor-abc", *rawResult.PaginationMeta.After)
}

func TestCollectedIncidentsEnvelope_LastPage(t *testing.T) {
	// The final page omits `after` entirely.
	body := []byte(`{
		"incidents": [{"id": "inc_3"}],
		"pagination_meta": {"page_size": 250}
	}`)
	rawResult := collectedIncidents{}
	require.NoError(t, json.Unmarshal(body, &rawResult))
	require.Len(t, rawResult.Incidents, 1)
	require.NotNil(t, rawResult.PaginationMeta)
	assert.Nil(t, rawResult.PaginationMeta.After)
}
