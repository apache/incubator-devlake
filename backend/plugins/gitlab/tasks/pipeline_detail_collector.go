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

var CollectApiPipelineDetailsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPipelines",
	EntryPoint:       CollectApiPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipeline data from gitlab api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiPipelineDetails(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	collectorWithState, err := helper.NewApiCollectorWithState(*rawDataSubTaskArgs, data.CreatedDateAfter)
	if err != nil {
		return err
	}

	incremental := collectorWithState.IsIncremental()

	iterator, err := GetPipelinesIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Input:              iterator,
		Incremental:        incremental,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines/{{ .Input.GitlabId }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("with_stats", "true")
			return query, nil
		},
		ResponseParser: GetOneRawMessageFromResponse,
		AfterResponse:  ignoreHTTPStatus403, // ignore 403 for CI/CD disable
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

func GetPipelinesIterator(taskCtx plugin.SubTaskContext, collectorWithState *helper.ApiCollectorStateManager) (*helper.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	clauses := []dal.Clause{
		dal.Select("gp.gitlab_id,gp.gitlab_id as iid"),
		dal.From("_tool_gitlab_pipelines gp"),
		dal.Where(
			`gp.project_id = ? and gp.connection_id = ?`,
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}
	if collectorWithState.CreatedDateAfter != nil {
		clauses = append(clauses, dal.Where("gitlab_created_at > ?", *collectorWithState.CreatedDateAfter))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}
