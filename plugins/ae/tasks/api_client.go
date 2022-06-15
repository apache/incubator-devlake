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
	"time"

	"github.com/apache/incubator-devlake/plugins/ae/api"
	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func CreateApiClient(taskCtx core.TaskContext, connection *models.AeConnection) (*helper.ApiAsyncClient, error) {
	// load and process cconfiguration
	endpoint := connection.Endpoint
	appId := connection.AppId
	secretKey := connection.SecretKey
	proxy := connection.Proxy

	apiClient, err := helper.NewApiClient(endpoint, nil, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetBeforeFunction(func(req *http.Request) error {
		nonceStr := core.RandLetterBytes(8)
		timestamp := fmt.Sprintf("%v", time.Now().Unix())
		sign := api.GetSign(req.URL.Query(), appId, secretKey, nonceStr, timestamp)
		req.Header.Set("x-ae-app-id", appId)
		req.Header.Set("x-ae-timestamp", timestamp)
		req.Header.Set("x-ae-nonce-str", nonceStr)
		req.Header.Set("x-ae-sign", sign)
		return nil
	})

	// create ae api client
	asyncApiCLient, err := helper.CreateAsyncApiClient(taskCtx, apiClient, nil)
	if err != nil {
		return nil, err
	}

	return asyncApiCLient, nil
}
