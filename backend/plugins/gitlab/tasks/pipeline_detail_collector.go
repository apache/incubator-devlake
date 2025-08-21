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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectApiPipelineDetailsMeta)
}

const RAW_PIPELINE_DETAILS_TABLE = "gitlab_api_pipeline_details"

var CollectApiPipelineDetailsMeta = plugin.SubTaskMeta{
	Name:             "Collect Pipeline Details",
	EntryPoint:       CollectApiPipelineDetails,
	EnabledByDefault: true,
	Description:      "Collect pipeline details data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiPipelinesMeta, &ExtractApiChildPipelinesMeta},
}

func CollectApiPipelineDetails(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_DETAILS_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	tickInterval, err := helper.CalcTickInterval(200, 1*time.Minute)
	if err != nil {
		return err
	}

	iterator, err := GetPipelinesIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		MinTickInterval:    &tickInterval,
		Input:              iterator,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines/{{ .Input.PipelineId }}",
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

func GetPipelinesIterator(taskCtx plugin.SubTaskContext, apiCollector *helper.StatefulApiCollector) (*helper.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	clauses := []dal.Clause{
		dal.Select("gp.pipeline_id"),
		dal.From("_tool_gitlab_pipeline_projects gp"),
		dal.Where(
			`gp.project_id = ? and gp.connection_id = ?`,
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("gitlab_updated_at > ?", *apiCollector.GetSince()))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(PipelineInput{}))
}

type PipelineInput struct {
	PipelineId int
}
