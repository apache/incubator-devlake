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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PULL_REQUEST_ACTIVITIES_TABLE = "bitbucket_server_api_pull_request_activities"

var CollectApiPrActivitiesMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequestsActivities",
	EntryPoint:       CollectApiPullRequestsActivities,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Collect pull requests activities data from Bitbucket Server api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	ProductTables:    []string{RAW_PULL_REQUEST_ACTIVITIES_TABLE},
}

func CollectApiPullRequestsActivities(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_ACTIVITIES_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	iterator, err := GetPullRequestsIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs:    *rawDataSubTaskArgs,
		ApiClient:             data.ApiClient,
		Incremental:           false,
		PageSize:              100,
		GetNextPageCustomData: GetNextPageCustomData,
		Query:                 GetQueryForNextPage,
		Input:                 iterator,
		UrlTemplate:           "rest/api/1.0/projects/{{ .Params.FullName }}/pull-requests/{{ .Input.BitbucketId }}/activities",
		ResponseParser:        GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
