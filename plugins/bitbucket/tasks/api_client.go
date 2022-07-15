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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func CreateApiClient(taskCtx core.TaskContext, connection *models.BitbucketConnection) (*helper.ApiAsyncClient, error) {
	// load configuration
	token := connection.GetEncodedToken()
	// create synchronize api client so we can calculate api rate limit dynamically
	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, nil, 0, connection.Proxy)
	if err != nil {
		return nil, err
	}
	// Rotates token on each request.
	apiClient.SetBeforeFunction(func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %v", token))
		return nil
	})
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth configuration")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimitPerHour,
		Method:               http.MethodGet,
		DynamicRateLimit: func(res *http.Response) (int, time.Duration, error) {

			return len(token), 1 * time.Hour, nil
		},
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
