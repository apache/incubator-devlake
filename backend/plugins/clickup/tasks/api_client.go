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
	"os"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	api "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

const CLICKUP_ACCESS_TOKEN_ENV = "CLICKUP_ACCESS_TOKEN"
const CLICKUP_TEAM_ID_ENV = "CLICKUP_TEAM_ID"

func NewClickupConnection() *models.ClickupConnection {
	accessToken := os.Getenv(CLICKUP_ACCESS_TOKEN_ENV)

	if accessToken == "" {
		panic(fmt.Sprintf("A personal clickup access token must be set in env var: %s", CLICKUP_ACCESS_TOKEN_ENV))
	}

	teamId := os.Getenv(CLICKUP_TEAM_ID_ENV)

	if teamId == "" {
		panic(fmt.Sprintf("Clickup team ID must be set in env var: %s", CLICKUP_TEAM_ID_ENV))
	}
	connection := models.ClickupConnection{
		RestConnection: api.RestConnection{
			Endpoint:         "https://api.clickup.com/api",
			Proxy:            "",
			RateLimitPerHour: 0,
		},
		TeamId: teamId,
		AccessToken: api.AccessToken{
			Token: accessToken,
		},
	}
	return &connection
}

func NewClickupApiClient(taskCtx plugin.TaskContext, connection *models.ClickupConnection) (*api.ApiAsyncClient, errors.Error) {
	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": connection.Token,
	}
	apiClient, err := api.NewApiClient(taskCtx.GetContext(), connection.Endpoint, headers, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) errors.Error {
		if res.StatusCode == http.StatusUnauthorized {
			return errors.HttpStatus(res.StatusCode).New("authentication failed, please check your AccessToken")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &api.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimitPerHour,
		DynamicRateLimit: func(res *http.Response) (int, time.Duration, errors.Error) {
			rateLimitHeader := res.Header.Get("RateLimit-Limit")
			if rateLimitHeader == "" {
				// use default
				return 0, 0, nil
			}
			rateLimit, err := strconv.Atoi(rateLimitHeader)
			if err != nil {
				return 0, 0, errors.Default.Wrap(err, "failed to parse RateLimit-Limit header")
			}
			// seems like {{ .plugin-ame }} rate limit is on minute basis
			return rateLimit, 1 * time.Minute, nil
		},
	}
	asyncApiClient, err := api.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}
	return asyncApiClient, nil
}
