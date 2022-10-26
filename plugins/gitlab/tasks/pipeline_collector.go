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
	goerror "errors"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
	"net/url"
)

const RAW_PIPELINE_TABLE = "gitlab_api_pipeline"

var CollectApiPipelinesMeta = core.SubTaskMeta{
	Name:             "collectApiPipelines",
	EntryPoint:       CollectApiPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipeline data from gitlab api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func CollectApiPipelines(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	if since == nil {
		var latestUpdated models.GitlabPipeline
		clause := []dal.Clause{
			dal.Orderby("gitlab_updated_at DESC"),
		}
		err := db.First(&latestUpdated, clause...)
		if err != nil && !goerror.Is(err, gorm.ErrRecordNotFound) {
			return errors.Default.Wrap(err, "failed to get latest gitlab pipeline record")
		}
		if latestUpdated.GitlabId > 0 {
			since = latestUpdated.GitlabUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Concurrency:        5,
		PageSize:           100,
		Incremental:        incremental,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if since != nil {
				query.Set("updated_after", since.String())
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

	return collector.Execute()
}
