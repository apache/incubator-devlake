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
	"net/http"
	"net/url"
)

var _ plugin.SubTaskEntryPoint = CollectAdditionalFilemetrics

const RAW_FILEMETRICS_ADDITIONAL_TABLE = "sonarqube_api_filemetrics_additional"

func CollectAdditionalFilemetrics(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FILEMETRICS_ADDITIONAL_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "measures/component_tree",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("component", data.Options.ProjectKey)
			query.Set("qualifiers", "FIL")
			query.Set("metricKeys", "duplicated_lines_density,duplicated_blocks,duplicated_lines, duplicated_files, complexity, cognitive_complexity, effort_to_reach_maintainability_rating_a, lines")
			query.Set("p", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("ps", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Data []json.RawMessage `json:"components"`
			}
			err := helper.UnmarshalResponse(res, &resData)
			return resData.Data, err
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectAdditionalFilemetricsMeta = plugin.SubTaskMeta{
	Name:             "CollectAdditionalFilemetrics",
	EntryPoint:       CollectAdditionalFilemetrics,
	EnabledByDefault: true,
	Description:      "Collect Filemetrics data from Sonarqube api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
