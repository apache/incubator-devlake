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
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

const RAW_BUILD_DEFINITION_TABLE = "azure_api_build_definitions"

var CollectApiBuildDefinitionMeta = plugin.SubTaskMeta{
	Name:        "collectApiBuild",
	EntryPoint:  CollectApiBuildDefinitions,
	Required:    true,
	Description: "Collect BuildDefinition data from Azure api",
	DomainTypes: []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiBuildDefinitions(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*AzureTaskData)

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AzureApiParams{
				ConnectionId: data.Options.ConnectionId,
				Project:      data.Options.Project,
			},
			Table: RAW_BUILD_DEFINITION_TABLE,
		},
		ApiClient: data.ApiClient,

		UrlTemplate: "{{ .Params.Project }}/_apis/build/definitions?api-version=7.1-preview.7",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Builds []json.RawMessage `json:"value"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.Builds, err
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
