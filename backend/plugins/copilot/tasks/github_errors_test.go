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
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildGitHubApiErrorExtractsMessageFromJson(t *testing.T) {
	err := buildGitHubApiError(http.StatusForbidden, "octodemo", []byte(`{"message":"nope"}`), "")
	require.Contains(t, err.Error(), "403")
	require.Contains(t, err.Error(), "Details: nope")
}

func TestBuildGitHubApiErrorIncludesRetryAfter(t *testing.T) {
	err := buildGitHubApiError(http.StatusTooManyRequests, "octodemo", []byte(""), "120")
	require.Contains(t, err.Error(), "429")
	require.Contains(t, err.Error(), "Retry-After: 120")
}

func TestBuildGitHubApiErrorStatusSpecificPrefixes(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   string
	}{
		{name: "403", status: http.StatusForbidden, want: "403 Forbidden"},
		{name: "404", status: http.StatusNotFound, want: "404 Not Found"},
		{name: "422", status: http.StatusUnprocessableEntity, want: "422 Unprocessable Entity"},
		{name: "429", status: http.StatusTooManyRequests, want: "429"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := buildGitHubApiError(tt.status, "octodemo", []byte("boom"), "")
			require.Contains(t, err.Error(), tt.want)
		})
	}
}
