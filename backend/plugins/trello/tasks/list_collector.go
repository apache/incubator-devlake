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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

const RAW_LIST_TABLE = "trello_lists"

var _ plugin.SubTaskEntryPoint = CollectList

var CollectListMeta = plugin.SubTaskMeta{
	Name:             "CollectList",
	EntryPoint:       CollectList,
	EnabledByDefault: true,
	Description:      "Collect list data from Trello api",
	DomainTypes:      []string{},
}

func CollectList(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TrelloTaskData)

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TrelloApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_LIST_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "1/boards/{{ .Params.BoardId }}/lists",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data []json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			return data, err
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
