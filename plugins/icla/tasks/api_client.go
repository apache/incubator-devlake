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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
)

const ENDPOINT = "https://people.apache.org/"

func NewIclaApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, errors.Error) {
	// load and process configuration
	token := taskCtx.GetConfig("ICLA_TOKEN")
	if token == "" {
		println("invalid ICLA_TOKEN, but ignore this error now")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("ICLA_API_REQUESTS_PER_HOUR"), 18000)
	if err != nil {
		return nil, err
	}
	proxy := taskCtx.GetConfig("ICLA_PROXY")

	// real request apiClient
	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), ENDPOINT, nil, 0, proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	// set token
	apiClient.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})

	// create async api client
	asyncApiClient, err := helper.CreateAsyncApiClient(taskCtx, apiClient, &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
	})
	if err != nil {
		return nil, err
	}

	return asyncApiClient, nil
}
