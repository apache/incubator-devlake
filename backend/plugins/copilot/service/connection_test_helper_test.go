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

package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildGitHubApiError_Forbidden(t *testing.T) {
	err := buildGitHubApiError(http.StatusForbidden, "octodemo", []byte(`{"message":"Resource not accessible"}`), "")
	assert.Equal(t, http.StatusForbidden, err.GetType().GetHttpCode())
	assert.Contains(t, err.Error(), "Forbidden")
	assert.Contains(t, err.Error(), "manage_billing:copilot")
}

func TestBuildGitHubApiError_NotFound(t *testing.T) {
	err := buildGitHubApiError(http.StatusNotFound, "octodemo", []byte(`{"message":"Not Found"}`), "")
	assert.Equal(t, http.StatusNotFound, err.GetType().GetHttpCode())
	assert.Contains(t, err.Error(), "octodemo")
}

func TestBuildGitHubApiError_UnprocessableEntity(t *testing.T) {
	err := buildGitHubApiError(http.StatusUnprocessableEntity, "octodemo", []byte(`{"message":"Metrics disabled"}`), "")
	assert.Equal(t, http.StatusUnprocessableEntity, err.GetType().GetHttpCode())
	assert.Contains(t, err.Error(), "Unprocessable")
	assert.Contains(t, err.Error(), "Metrics disabled")
}

func TestBuildGitHubApiError_TooManyRequests(t *testing.T) {
	err := buildGitHubApiError(http.StatusTooManyRequests, "octodemo", []byte(`{"message":"Slow down"}`), "60")
	assert.Equal(t, http.StatusTooManyRequests, err.GetType().GetHttpCode())
	assert.Contains(t, err.Error(), "429")
	assert.Contains(t, err.Error(), "60 seconds")
}
