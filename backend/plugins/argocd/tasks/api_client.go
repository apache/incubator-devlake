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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

func NewArgocdApiClient(taskCtx plugin.TaskContext, connection *models.ArgocdConnection) (*api.ApiAsyncClient, errors.Error) {
	apiClient, err := api.NewApiClientFromConnection(
		taskCtx.GetContext(),
		taskCtx,
		connection,
	)
	if err != nil {
		return nil, err
	}

	apiClient.SetData("Token", connection.Token)

	// use standard creator to ensure scheduler initialized
	asyncClient, err2 := api.CreateAsyncApiClient(taskCtx, apiClient, nil)
	if err2 != nil {
		return nil, err2
	}
	return asyncClient, nil
}
