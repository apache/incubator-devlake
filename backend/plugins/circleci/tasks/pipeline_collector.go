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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PIPELINE_TABLE = "circleci_api_pipelines"

var _ plugin.SubTaskEntryPoint = CollectPipelines

var CollectPipelinesMeta = plugin.SubTaskMeta{
	Name:             "collectPipelines",
	EntryPoint:       CollectPipelines,
	EnabledByDefault: true,
	Description:      "collect circleci pipelines",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	logger := taskCtx.GetLogger()
	timeAfter := rawDataSubTaskArgs.Ctx.TaskContext().SyncPolicy().TimeAfter
	logger.Info("collect pipelines")
	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			PageSize:              int(data.Options.PageSize),
			GetNextPageCustomData: ExtractNextPageToken,
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "/v2/project/{{ .Params.ProjectSlug }}/pipeline",
				Query:       BuildQueryParamsWithPageToken,
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					data := CircleciPageTokenResp[[]json.RawMessage]{}
					err := api.UnmarshalResponse(res, &data)

					if err != nil {
						return nil, err
					}
					filteredItems := []json.RawMessage{}
					for _, item := range data.Items {
						pipelineCreatedAt, err := extractCreatedAt(item)

						if err != nil {
							return nil, err
						}
						if pipelineCreatedAt.Before(*timeAfter) {
							return filteredItems, api.ErrFinishCollect
						}
						filteredItems = append(filteredItems, item)
					}
					return filteredItems, nil
				},
			},
			GetCreated: extractCreatedAt,
		},
	})
	if err != nil {
		logger.Error(err, "collect pipelines error")
		return err
	}
	return collector.Execute()
}
