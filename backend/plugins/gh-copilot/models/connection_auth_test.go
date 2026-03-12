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

package models

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGhCopilotConn_SetupAuthentication_BearerPrefix(t *testing.T) {
	conn := &GhCopilotConn{Token: "Bearer abc"}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err2 := conn.SetupAuthentication(req)
	require.NoError(t, err2)
	require.Equal(t, "Bearer abc", req.Header.Get("Authorization"))
}

func TestGhCopilotConn_SetupAuthentication_TokenPrefix(t *testing.T) {
	conn := &GhCopilotConn{Token: "token abc"}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err2 := conn.SetupAuthentication(req)
	require.NoError(t, err2)
	require.Equal(t, "token abc", req.Header.Get("Authorization"))
}

func TestGhCopilotConn_SetupAuthentication_RawToken(t *testing.T) {
	conn := &GhCopilotConn{Token: "abc"}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err2 := conn.SetupAuthentication(req)
	require.NoError(t, err2)
	require.Equal(t, "Bearer abc", req.Header.Get("Authorization"))
}

func TestGhCopilotConn_SetupAuthentication_TrimsWhitespace(t *testing.T) {
	conn := &GhCopilotConn{Token: "  abc  "}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	err2 := conn.SetupAuthentication(req)
	require.NoError(t, err2)
	require.Equal(t, "Bearer abc", req.Header.Get("Authorization"))
}
