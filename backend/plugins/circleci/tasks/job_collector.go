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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_JOB_TABLE = "circleci_api_jobs"

var _ plugin.SubTaskEntryPoint = CollectJobs

var CollectJobsMeta = plugin.SubTaskMeta{
	Name:             "collectJobs",
	EntryPoint:       CollectJobs,
	EnabledByDefault: true,
	Description:      "collect circleci jobs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect jobs")

	clauses := []dal.Clause{
		dal.Select("id, pipeline_id"),
		dal.From(&models.CircleciWorkflow{}),
		dal.Where("_tool_circleci_workflows.connection_id = ? and _tool_circleci_workflows.project_slug = ? ", data.Options.ConnectionId, data.Options.ProjectSlug),
	}

	db := taskCtx.GetDal()
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.CircleciWorkflow{}))
	if err != nil {
		return err
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "/v2/workflow/{{ .Input.Id }}/job",
		Input:              iterator,
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			res := CircleciPageTokenResp[any]{}
			err := api.UnmarshalResponse(prevPageResponse, &res)
			if err != nil {
				return nil, err
			}
			if res.NextPageToken == "" {
				return nil, api.ErrFinishCollect
			}
			return res.NextPageToken, nil
		},
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
				query.Set("page_token", reqData.CustomData.(string))
			}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			data := CircleciPageTokenResp[[]json.RawMessage]{}
			err := api.UnmarshalResponse(res, &data)
			return data.Items, err
		},
	})
	if err != nil {
		logger.Error(err, "collect jobs error")
		return err
	}
	return collector.Execute()
}
