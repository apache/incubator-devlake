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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_COMMITS_TABLE = "bitbucket_server_api_commits"

var CollectApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectApiCommits",
	EntryPoint:       CollectApiCommits,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Collect commits data from Bitbucket Server api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	ProductTables:    []string{RAW_COMMITS_TABLE},
}

func CollectApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	iterator, err := GetBranchesIterator(taskCtx)
	if err != nil {
		return err
	}
	defer iterator.Close()

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMITS_TABLE)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs:    *rawDataSubTaskArgs,
		ApiClient:             data.ApiClient,
		Incremental:           false,
		PageSize:              100,
		GetNextPageCustomData: GetNextPageCustomData,
		Query:                 GetQueryForNextPage,
		Input:                 iterator,
		UrlTemplate:           "rest/api/1.0/projects/{{ .Params.FullName }}/commits?until={{ .Input.Branch }}",
		ResponseParser:        GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
