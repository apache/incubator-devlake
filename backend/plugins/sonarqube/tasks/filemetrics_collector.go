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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_FILEMETRICS_TABLE = "sonarqube_api_filemetrics"

var _ plugin.SubTaskEntryPoint = CollectFilemetrics

func CollectFilemetrics(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	logger.Info("collect filemetrics")

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FILEMETRICS_TABLE)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		UrlTemplate:        "measures/component_tree",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("component", data.Options.ProjectKey)
			query.Set("qualifiers", "FIL")
			query.Set("metricKeys", "code_smells,sqale_index,sqale_rating,bugs,reliability_rating,vulnerabilities,security_rating,security_hotspots,security_hotspots_reviewed,security_review_rating,ncloc,coverage,uncovered_lines,lines_to_cover")
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

var CollectFilemetricsMeta = plugin.SubTaskMeta{
	Name:             "CollectFilemetrics",
	EntryPoint:       CollectFilemetrics,
	EnabledByDefault: true,
	Description:      "Collect Filemetrics data from Sonarqube api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
