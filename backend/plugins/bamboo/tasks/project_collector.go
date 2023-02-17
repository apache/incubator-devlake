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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"io"
	"net/http"
	"net/url"
)

const RAW_PROJECT_TABLE = "bamboo_project"

var _ plugin.SubTaskEntryPoint = CollectProject

func CollectProject(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	collectorWithState, err := helper.NewApiCollectorWithState(*rawDataSubTaskArgs, nil)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		Incremental: incremental,
		ApiClient:   data.ApiClient,
		// TODO write which api would you want request
		UrlTemplate: "project/{{ .Params.ProjectKey }}.json",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("showEmpty", fmt.Sprintf("%v", true))
			return query, nil
		},
		GetTotalPages: func(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
			return 1, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			res.Body.Close()
			return []json.RawMessage{body}, nil
		},
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectProjectMeta = plugin.SubTaskMeta{
	Name:             "CollectProject",
	EntryPoint:       CollectProject,
	EnabledByDefault: true,
	Description:      "Collect Project data from Bamboo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
