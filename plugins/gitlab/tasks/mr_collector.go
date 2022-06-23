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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_MERGE_REQUEST_TABLE = "gitlab_api_merge_requests"

var CollectApiMergeRequestsMeta = core.SubTaskMeta{
	Name:             "collectApiMergeRequests",
	EntryPoint:       CollectApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Collect merge requests data from gitlab api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func CollectApiMergeRequests(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/merge_requests",
		Query:              GetQuery,
		GetTotalPages:      GetTotalPagesFromResponse,
		ResponseParser:     GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
