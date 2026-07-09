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
	stdctx "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	corectx "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
)

// TestConnectionResult represents the payload returned by the connection test endpoints.
type TestConnectionResult struct {
	Success                 bool   `json:"success"`
	Message                 string `json:"message"`
	MemberCount             int    `json:"memberCount,omitempty"`
	HasEnterpriseAnalytics  bool   `json:"hasEnterpriseAnalytics,omitempty"`
}

// TestConnection exercises the Cursor Admin API to validate credentials.
func TestConnection(ctx stdctx.Context, br corectx.BasicRes, connection *models.CursorConnection) (*TestConnectionResult, errors.Error) {
	if connection == nil {
		return nil, errors.BadInput.New("connection is required")
	}

	connection.Normalize()
	if strings.TrimSpace(connection.Token) == "" {
		return nil, errors.BadInput.New("token is required")
	}

	apiClient, err := helper.NewApiClientFromConnection(ctx, br, connection)
	if err != nil {
		return nil, err
	}

	if userKeyErr := detectUserApiKey(apiClient); userKeyErr != nil {
		return &TestConnectionResult{
			Success: false,
			Message: userKeyErr.Error(),
		}, nil
	}

	res, err := apiClient.Get("teams/members", nil, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to reach Cursor Admin API")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return &TestConnectionResult{
			Success: false,
			Message: "Invalid Team API key. Create a Team Admin key in Dashboard → API Keys (admin:* scope), not a User API key from Settings → Integrations.",
		}, nil
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return &TestConnectionResult{
			Success: false,
			Message: fmt.Sprintf("unexpected status %d: %s", res.StatusCode, string(body)),
		}, nil
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, errors.Default.Wrap(readErr, "failed to read response body")
	}

	var response struct {
		TeamMembers []json.RawMessage `json:"teamMembers"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return &TestConnectionResult{
			Success: false,
			Message: fmt.Sprintf("failed to parse response: %v", jsonErr),
		}, nil
	}

	hasEnterpriseAnalytics := probeEnterpriseAnalytics(apiClient)

	return &TestConnectionResult{
		Success:                true,
		Message:                "Connection validated successfully",
		MemberCount:            len(response.TeamMembers),
		HasEnterpriseAnalytics: hasEnterpriseAnalytics,
	}, nil
}

func detectUserApiKey(apiClient *helper.ApiClient) errors.Error {
	res, err := apiClient.Get("v1/me", nil, nil)
	if err != nil || res == nil {
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return errors.BadInput.New("this is a User API key. Create a Team Admin key in Dashboard → API Keys (admin:* scope)")
	}
	return nil
}

func probeEnterpriseAnalytics(apiClient *helper.ApiClient) bool {
	res, err := apiClient.Get("analytics/team/dau?startDate=7d&endDate=today", nil, nil)
	if err != nil || res == nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}
