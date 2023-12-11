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
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PULL_REQUEST_TABLE = "bitbucket_server_api_pull_requests"

// this struct should be moved to `bitbucket_api_common.go`

var CollectApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequests",
	EntryPoint:       CollectApiPullRequests,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Collect PullRequests data from Bitbucket Server api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient: data.ApiClient,
		PageSize:  25,
		GetNextPageCustomData: func(prevReqData *helper.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			var rawMessages struct {
				NextPageStart int  `json:"nextPageStart"`
				IsLastPage    bool `json:"isLastPage"`
			}
			err := decodeResponse(prevPageResponse, &rawMessages)
			if err != nil {
				return nil, err
			}

			if rawMessages.IsLastPage {
				return nil, api.ErrFinishCollect
			}

			return strconv.Itoa(rawMessages.NextPageStart), nil
		},
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("pagelen", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("sort", "created_on")

			if reqData.CustomData != nil {
				query.Set("start", reqData.CustomData.(string))
			}
			return query, nil
		},
		UrlTemplate:    "rest/api/1.0/projects/{{ .Params.FullName }}/pull-requests",
		ResponseParser: GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
