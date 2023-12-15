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
	"net/url"
)

const RAW_USER_TABLE = "circleci_api_users"

var _ plugin.SubTaskEntryPoint = CollectAccounts

var CollectAccountsMeta = plugin.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "collect circleci accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect accounts")
	collectorWithState, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	if err := collectorWithState.InitCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "/v2/me/collaborations",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data []json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			return data, err
		},
	}); err != nil {
		logger.Error(err, "collect user error")
		return err
	}
	return collectorWithState.Execute()
}
