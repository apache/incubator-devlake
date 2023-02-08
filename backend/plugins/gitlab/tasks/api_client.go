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
	gocontext "context"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

// NewApiClientFromConnectionWithTest creates ApiClient based on given connection.
// and then it will test and try to correct the connection
// The connection must
func NewApiClientFromConnectionWithTest(
	ctx gocontext.Context,
	br context.BasicRes,
	connection *models.GitlabConn,
) (*api.ApiClient, errors.Error) {
	connection.IsPrivateToken = false
	// test connection
	apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), br, connection)
	if err != nil {
		return nil, errors.Convert(err)
	}

	resBody := &models.ApiUserResponse{}
	res, err := apiClient.Get("user", nil, nil)
	if res.StatusCode != http.StatusUnauthorized {
		if err != nil {
			return nil, errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, resBody)
		if err != nil {
			return nil, errors.Convert(err)
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
		}
	} else {
		connection.IsPrivateToken = true
		res, err = apiClient.Get("user", nil, nil)
		if err != nil {
			return nil, errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, resBody)
		if err != nil {
			return nil, errors.Convert(err)
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection[PrivateToken]")
		}
	}
	connection.UserId = resBody.Id

	return apiClient, nil
}

func CreateGitlabAsyncApiClient(
	taskCtx plugin.TaskContext,
	apiClient *api.ApiClient,
	connection *models.GitlabConnection,
) (*api.ApiAsyncClient, errors.Error) {
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
			// seems like gitlab rate limit is on minute basis
			if rateLimit > 200 {
				return 200, 1 * time.Minute, nil
			} else {
				return rateLimit, 1 * time.Minute, nil
			}
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

func NewGitlabApiClient(taskCtx plugin.TaskContext, connection *models.GitlabConnection) (*api.ApiAsyncClient, errors.Error) {
	apiClient, err := NewApiClientFromConnectionWithTest(taskCtx.GetContext(), taskCtx, &connection.GitlabConn)
	if err != nil {
		return nil, err
	}

	return CreateGitlabAsyncApiClient(taskCtx, apiClient, connection)
}

func ignoreHTTPStatus403(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusForbidden {
		return api.ErrIgnoreAndContinue
	}
	return nil
}
