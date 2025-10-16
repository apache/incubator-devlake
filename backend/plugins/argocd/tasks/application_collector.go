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
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

const RAW_APPLICATION_TABLE = "argocd_api_applications"

var _ plugin.SubTaskEntryPoint = CollectApplications

var CollectApplicationsMeta = plugin.SubTaskMeta{
	Name:             "collectApplications",
	EntryPoint:       CollectApplications,
	EnabledByDefault: true,
	Description:      "Collect applications from ArgoCD API",
}

func CollectApplications(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ArgocdTaskData)

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_APPLICATION_TABLE,
			Params: models.ArgocdApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.ApplicationName,
			},
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "/applications/{{.Params.Name}}",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("refresh", "false")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var application json.RawMessage
			err := api.UnmarshalResponse(res, &application)
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{application}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
