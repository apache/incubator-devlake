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
	"log"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argo/models"
)

const RAW_DEPLOYMENT_TABLE = "argo_api_deployments"

var CollectDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "collect_deployments",
	EntryPoint:       CollectApiDeployments,
	EnabledByDefault: true,
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_DEPLOYMENT_TABLE},
}

// Coletor principal
func CollectApiDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	log.Println("[ARGO] Iniciando plugin de collect.")

	data := taskCtx.GetData().(*models.ArgoTaskData)
	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: models.ArgoApiParams{
			ConnectionId: data.Options.ConnectionId,
			ProjectId:    data.Options.ProjectId,
		},
		Table: RAW_DEPLOYMENT_TABLE,
	})
	if err != nil {
		return err
	}

	log.Println("[ARGO] ConnectionId: " + fmt.Sprintf("%d", data.Options.ConnectionId))
	log.Println("[ARGO] ProjectId: " + data.Options.ProjectId)

	err = apiCollector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "/api/v1/workflows/argo/{{ .Params.ProjectId }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("offset", fmt.Sprintf("%v", reqData.Pager.Page*reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages:  models.GetTotalPagesFromResponse,
		ResponseParser: models.GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	log.Println("[ARGO] Finalizado plugin de collect.")

	return apiCollector.Execute()
}
