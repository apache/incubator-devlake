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
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

const RAW_PROJECT_METRICS_HISTORY_TABLE = "sonarqube_api_project_metrics_history"

const metricsToCollect = "coverage,ncloc,bugs,vulnerabilities,code_smells,security_hotspots," +
	"duplicated_lines_density,sqale_rating,reliability_rating,security_rating,complexity,cognitive_complexity"

var _ plugin.SubTaskEntryPoint = CollectProjectMetricsHistory

func CollectProjectMetricsHistory(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	logger.Info("collect project metrics history")

	data := taskCtx.GetData().(*SonarqubeTaskData)
	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: models.SonarqubeApiParams{
			ConnectionId: data.Options.ConnectionId,
			ProjectKey:   data.Options.ProjectKey,
		},
		Table: RAW_PROJECT_METRICS_HISTORY_TABLE,
	})
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    1000,
		UrlTemplate: "measures/search_history",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("component", data.Options.ProjectKey)
			query.Set("metrics", metricsToCollect)
			query.Set("ps", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("p", fmt.Sprintf("%v", reqData.Pager.Page))
			if apiCollector.GetSince() != nil {
				query.Set("from", apiCollector.GetSince().UTC().Format("2006-01-02"))
			}
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var body json.RawMessage
			err := helper.UnmarshalResponse(res, &body)
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}

var CollectProjectMetricsHistoryMeta = plugin.SubTaskMeta{
	Name:             "CollectProjectMetricsHistory",
	EntryPoint:       CollectProjectMetricsHistory,
	EnabledByDefault: true,
	Description:      "Collect project-level metric history from SonarQube measures/search_history API",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
