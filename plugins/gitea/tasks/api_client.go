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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"time"
)

func NewGiteaApiClient(taskCtx core.TaskContext, connection *models.GiteaConnection) (*helper.ApiAsyncClient, error) {

	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, nil, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}

	apiClient.SetBeforeFunction(func(req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", connection.Token))
		return nil
	})

	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		nil,
	)
	apiClient.SetTimeout(30 * time.Second)
	if err != nil {
		return nil, err
	}
	return asyncApiClient, nil
}
