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
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_MERGE_REQUEST_TABLE = "gitlab_api_merge_requests"

func init() {
	RegisterSubtaskMeta(&CollectApiMergeRequestsMeta)
}

var CollectApiMergeRequestsMeta = plugin.SubTaskMeta{
	Name:             "Collect Merge Requests",
	EntryPoint:       CollectApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Collect merge requests data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiIssuesMeta},
}

func CollectApiMergeRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	apiCollector, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:      data.ApiClient,
		PageSize:       100,
		UrlTemplate:    "projects/{{ .Params.ProjectId }}/merge_requests",
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query, err := GetQuery(reqData)
			if err != nil {
				return nil, err
			}
			if apiCollector.GetSince() != nil {
				query.Set("updated_after", apiCollector.GetSince().Format(time.RFC3339))
			}
			return query, nil
		},
	})

	if err != nil {
		return err
	}

	return apiCollector.Execute()
}
