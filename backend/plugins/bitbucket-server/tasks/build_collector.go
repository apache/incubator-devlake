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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_BUILDS_TABLE = "bitbucket_server_api_builds"

var CollectApiBuildsMeta = plugin.SubTaskMeta{
	Name:             "collectApiBuilds",
	EntryPoint:       CollectApiBuilds,
	EnabledByDefault: true,
	Description:      "Collect Builds data from Bitbucket Server api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUILDS_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	iterator, err := GetCommitsIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Input:       iterator,
		UrlTemplate: "rest/build-status/1.0/commits/stats/{{ .Input.CommitSha }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var rawMessage json.RawMessage
			err := decodeResponse(res, &rawMessage)
			if err != nil {
				return nil, err
			}

			return []json.RawMessage{rawMessage}, nil
		},
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
