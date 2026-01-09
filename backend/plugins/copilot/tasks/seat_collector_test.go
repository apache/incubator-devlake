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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCopilotSeatsFromResponse_WrappedObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://api.github.com/orgs/octodemo/copilot/billing/seats?page=1", nil)
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"total_seats":2,"seats":[{"assignee":{"login":"a"}},{"assignee":{"login":"b"}}]}`)),
		Request:    req,
	}

	msgs, err := parseCopilotSeatsFromResponse(res)
	require.NoError(t, err)
	require.Len(t, msgs, 2)
}

func TestParseCopilotSeatsFromResponse_Array(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://api.github.com/orgs/octodemo/copilot/billing/seats?page=1", nil)
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`[{"assignee":{"login":"a"}}]`)),
		Request:    req,
	}

	msgs, err := parseCopilotSeatsFromResponse(res)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
}
