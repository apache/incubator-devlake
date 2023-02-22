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
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_DEPLOY_TABLE = "bamboo_api_deploy"

var _ plugin.SubTaskEntryPoint = CollectDeploy

func CollectDeploy(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOY_TABLE)

	collectorWithState, err := helper.NewApiCollectorWithState(*rawDataSubTaskArgs, nil)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "deploy/project/all.json",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("showEmpty", fmt.Sprintf("%v", true))
			query.Set("max-result", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("start-index", fmt.Sprintf("%v", reqData.Pager.Skip))
			return query, nil
		},

		ResponseParser: helper.GetRawMessageArrayFromResponse,
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectDeployMeta = plugin.SubTaskMeta{
	Name:             "CollectDeploy",
	EntryPoint:       CollectDeploy,
	EnabledByDefault: true,
	Description:      "Collect Deploy data from Bamboo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
