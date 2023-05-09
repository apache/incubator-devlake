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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"golang.org/x/mod/semver"
)

const RAW_USER_TABLE = "gitlab_api_users"

var CollectAccountsMeta = plugin.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "collect gitlab users, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect gitlab users")

	// it means we can not use /members/all to get the data
	urlTemplate := "/projects/{{ .Params.ProjectId }}/members/all"
	if semver.Compare(data.ApiClient.GetData(models.GitlabApiClientData_ApiVersion).(string), "v13.11") < 0 {
		urlTemplate = "/projects/{{ .Params.ProjectId }}/members/"
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        urlTemplate,
		PageSize:           100,
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var items []json.RawMessage
			err := api.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})

	if err != nil {
		logger.Error(err, "collect user error")
		return err
	}

	return collector.Execute()
}
