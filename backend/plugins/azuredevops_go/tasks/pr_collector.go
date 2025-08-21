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
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectApiPullRequestsMeta)
}

const RawPullRequestTable = "azuredevops_go_api_pull_requests"

var CollectApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequests",
	EntryPoint:       CollectApiPullRequests,
	EnabledByDefault: true,
	Description:      "Collect PullRequests data from Azure DevOps API, supports timeFilter but not diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{},
	ProductTables:    []string{RawPullRequestTable},
}

func CollectApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*AzuredevopsTaskData)

	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Table:   RawPullRequestTable,
		Options: data.Options,
	}

	apiCollector, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/git/repositories/{{ .Params.RepositoryId }}/pullrequests?api-version=7.1",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("searchCriteria.status", "all")
			query.Set("$skip", fmt.Sprint(reqData.Pager.Skip))
			query.Set("$top", fmt.Sprint(reqData.Pager.Size))

			if apiCollector.GetSince() != nil {
				query.Set("searchCriteria.queryTimeRangeType", "created")
				query.Set("searchCriteria.minTime", apiCollector.GetSince().Format(time.RFC3339))
			}
			return query, nil
		},
		ResponseParser: ParseRawMessageFromValue,
		AfterResponse:  change203To401,
	})

	if err != nil {
		return err
	}

	return apiCollector.Execute()
}
