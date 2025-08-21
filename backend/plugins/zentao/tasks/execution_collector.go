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
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_EXECUTION_TABLE = "zentao_api_executions"

var _ plugin.SubTaskEntryPoint = CollectExecutions

func CollectExecutions(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	cursor, iterator, err := getExecutionIterator(taskCtx)
	if err != nil {
		return err
	}
	defer cursor.Close()
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_EXECUTION_TABLE,
		},
		Input:       iterator,
		ApiClient:   data.ApiClient,
		UrlTemplate: "/executions/{{ .Input.Id }}",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			if errors.Is(err, api.ErrEmptyResponse) {
				return nil, nil
			}
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{data}, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectExecutionMeta = plugin.SubTaskMeta{
	Name:             "collectExecutions",
	EntryPoint:       CollectExecutions,
	EnabledByDefault: true,
	Description:      "Collect Execution data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
