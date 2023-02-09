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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_ACCOUNTS_TABLE = "sonarqube_accounts"

var _ plugin.SubTaskEntryPoint = CollectAccounts

func CollectAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ACCOUNTS_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect accounts")

	collectorWithState, err := helper.NewApiCollectorWithState(*rawDataSubTaskArgs, data.CreatedDateAfter)
	if err != nil {
		return err
	}
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "users/search",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("p", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("ps", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Data []json.RawMessage `json:"users"`
			}
			err = helper.UnmarshalResponse(res, &resData)
			return resData.Data, err
		},
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

var CollectAccountsMeta = plugin.SubTaskMeta{
	Name:             "CollectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "Collect Accounts data from Sonarqube user api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
