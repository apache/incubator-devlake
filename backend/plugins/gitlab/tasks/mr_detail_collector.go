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
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_MERGE_REQUEST_DETAIL_TABLE = "gitlab_api_merge_request_details"

var CollectApiMergeRequestDetailsMeta = plugin.SubTaskMeta{
	Name:             "collectApiMergeRequestDetails",
	EntryPoint:       CollectApiMergeRequestDetails,
	EnabledByDefault: true,
	Description:      "Collect merge request Details data from gitlab api, supports timeFilter but not diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiMergeRequestDetails(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_DETAIL_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}

	iterator, err := GetMergeRequestDetailsIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		Incremental: false,
		Input:       iterator,
		UrlTemplate: "projects/{{ .Params.ProjectId }}/merge_requests/{{ .Input.Iid }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("with_stats", "true")
			return query, nil
		},
		ResponseParser: GetOneRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

func GetMergeRequestDetailsIterator(taskCtx plugin.SubTaskContext, collectorWithState *helper.ApiCollectorStateManager) (*helper.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	clauses := []dal.Clause{
		dal.Select("gmr.gitlab_id, gmr.iid"),
		dal.From("_tool_gitlab_merge_requests gmr"),
		dal.Where(
			`gmr.project_id = ? and gmr.connection_id = ? and gmr.is_detail_required = ?`,
			data.Options.ProjectId, data.Options.ConnectionId, true,
		),
	}
	if collectorWithState.LatestState.LatestSuccessStart != nil {
		clauses = append(clauses, dal.Where("gitlab_updated_at > ?", *collectorWithState.LatestState.LatestSuccessStart))
	} else if collectorWithState.TimeAfter != nil {
		clauses = append(clauses, dal.Where("gitlab_updated_at > ?", *collectorWithState.TimeAfter))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}
