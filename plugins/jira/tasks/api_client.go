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
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func NewJiraApiClient(taskCtx core.TaskContext, connection *models.JiraConnection) (*helper.ApiAsyncClient, error) {
	// decrypt connection first
	err := helper.DecryptConnection(connection)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt Auth AccessToken: %w", err)
	}

	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %v", connection.GetEncodedToken()),
	}
	apiClient, err := helper.NewApiClient(connection.Endpoint, headers, 0, connection.Proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your AccessToken")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimit,
	}
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	return asyncApiClient, nil
}

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

func GetJiraServerInfo(client *helper.ApiAsyncClient) (*models.JiraServerInfo, int, error) {
	res, err := client.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, 0, err
	}
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, res.StatusCode, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	serverInfo := &models.JiraServerInfo{}
	err = helper.UnmarshalResponse(res, serverInfo)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return serverInfo, res.StatusCode, nil
}
