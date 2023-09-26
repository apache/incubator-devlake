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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_FOLDERLESS_TABLE = "clickup_folderless_list"

var _ plugin.SubTaskEntryPoint = CollectIssue

func CollectFolderlessList(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FOLDERLESS_TABLE)

	collectorWithState, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs, data.CreatedDateAfter)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		Incremental: incremental,
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "v2/space/{{ .Params.SpaceId }}/list",
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			var res struct {
				LastPage bool `json:"last_page"`
			}
			err := api.UnmarshalResponse(prevPageResponse, &res)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", reqData.Pager.Page-1))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := struct {
				Lists[]json.RawMessage
			}{}
			err := helper.UnmarshalResponse(res, &body)
			if err != nil {
				return nil, err
			}
			return body.Lists, nil
		},
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectFolderlessListMeta = plugin.SubTaskMeta{
	Name:             "CollectFolderlessList",
	EntryPoint:       CollectFolderlessList,
	EnabledByDefault: true,
	Description:      "Collect FolderlessList data from Clickup api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
