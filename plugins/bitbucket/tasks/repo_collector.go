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
	"github.com/apache/incubator-devlake/plugins/helper"
	"io"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
)

const RAW_REPOSITORIES_TABLE = "bitbucket_api_repositories"

var CollectApiRepoMeta = core.SubTaskMeta{
	Name:        "collectApiRepo",
	EntryPoint:  CollectApiRepositories,
	Required:    true,
	Description: "Collect repositories data from Bitbucket api",
	DomainTypes: []string{core.DOMAIN_TYPE_CODE},
}

func CollectApiRepositories(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_REPOSITORIES_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,

		UrlTemplate: "repositories/{{ .Params.Owner }}/{{ .Params.Repo }}",
		Query:       GetQuery,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
