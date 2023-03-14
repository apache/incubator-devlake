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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PIPELINE_TABLE = "gitlab_api_pipeline"

var CollectApiPipelinesMeta = plugin.SubTaskMeta{
	Name:             "collectApiPipelines",
	EntryPoint:       CollectApiPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipeline data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}

	tickInterval, err := helper.CalcTickInterval(200, 1*time.Minute)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		MinTickInterval:    &tickInterval,
		PageSize:           100,
		Incremental:        incremental,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if collectorWithState.TimeAfter != nil {
				query.Set("updated_after", collectorWithState.TimeAfter.Format(time.RFC3339))
			}
			if incremental {
				query.Set("updated_after", collectorWithState.LatestState.LatestSuccessStart.Format(time.RFC3339))
			}
			query.Set("with_stats", "true")
			query.Set("sort", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: GetRawMessageFromResponse,
		AfterResponse:  ignoreHTTPStatus403, // ignore 403 for CI/CD disable
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
