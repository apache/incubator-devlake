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

const RAW_USERS_TABLE = "opsgenie_users"

type (
	collectedUsers struct {
		TotalCount int               `json:"totalCount"`
		Data       []json.RawMessage `json:"data"`
	}
)

var CollectUsersMeta = plugin.SubTaskMeta{
	Name:             "collectUsers",
	EntryPoint:       CollectUsers,
	EnabledByDefault: true,
	Description:      "collect Opsgenie users.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectUsers(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_USERS_TABLE,
		},
		ApiClient:   data.Client,
		UrlTemplate: "v2/users",
		PageSize:    100,
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}

			query.Set("sort", "createdAt")
			query.Set("order", "desc")
			query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
			query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Skip))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			rawResult := collectedUsers{}
			err := api.UnmarshalResponse(res, &rawResult)

			return rawResult.Data, err
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
