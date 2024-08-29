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
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
)

var accountIdGen *didgen.DomainIdGenerator
var projectIdGen *didgen.DomainIdGenerator
var pipelineIdGen *didgen.DomainIdGenerator
var jobIdGen *didgen.DomainIdGenerator

func getAccountIdGen() *didgen.DomainIdGenerator {
	if accountIdGen == nil {
		accountIdGen = didgen.NewDomainIdGenerator(&models.CircleciAccount{})
	}
	return accountIdGen
}

func getProjectIdGen() *didgen.DomainIdGenerator {
	if projectIdGen == nil {
		projectIdGen = didgen.NewDomainIdGenerator(&models.CircleciProject{})
	}
	return projectIdGen
}

func getWorkflowIdGen() *didgen.DomainIdGenerator {
	if pipelineIdGen == nil {
		pipelineIdGen = didgen.NewDomainIdGenerator(&models.CircleciWorkflow{})
	}
	return pipelineIdGen
}

func getJobIdGen() *didgen.DomainIdGenerator {
	if jobIdGen == nil {
		jobIdGen = didgen.NewDomainIdGenerator(&models.CircleciJob{})
	}
	return jobIdGen
}

type CircleciPageTokenResp[T any] struct {
	Items         T      `json:"items"`
	NextPageToken string `json:"next_page_token"`
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *CircleciTaskData) {
	data := taskCtx.GetData().(*CircleciTaskData)
	filteredData := *data
	filteredData.Options = &CircleciOptions{}
	*filteredData.Options = *data.Options
	params := models.CircleciApiParams{
		ConnectionId: data.Options.ConnectionId,
		ProjectSlug:  data.Options.ProjectSlug,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}

func findPipelineById(db dal.Dal, id string) (*models.CircleciPipeline, errors.Error) {
	if id == "" {
		return nil, errors.Default.New("id must not empty")
	}
	pipeline := &models.CircleciPipeline{}
	if err := db.First(pipeline, dal.Where("id = ?", id)); err != nil {
		return nil, err
	}
	return pipeline, nil
}

func ExtractNextPageToken(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	res := CircleciPageTokenResp[any]{}
	err := api.UnmarshalResponse(prevPageResponse, &res)
	if err != nil {
		return nil, err
	}
	if res.NextPageToken == "" {
		return nil, api.ErrFinishCollect
	}
	return res.NextPageToken, nil
}

func BuildQueryParamsWithPageToken(reqData *api.RequestData, _ *time.Time) (url.Values, errors.Error) {
	query := url.Values{}
	if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
		query.Set("page-token", pageToken)
	}
	return query, nil
}

func ParseCircleciPageTokenResp(res *http.Response) ([]json.RawMessage, errors.Error) {
	data := CircleciPageTokenResp[[]json.RawMessage]{}
	err := api.UnmarshalResponse(res, &data)
	return data.Items, err
}

func ignoreDeletedBuilds(res *http.Response) errors.Error {
	// CircleCI API will return a 404 response for a workflow/job that has been deleted
	// due to their data retention policy. We should ignore these errors.
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}

func extractCreatedAt(item json.RawMessage) (time.Time, errors.Error) {
	var entity struct {
		CreatedAt time.Time `json:"created_at"`
	}
	if err := json.Unmarshal(item, &entity); err != nil {
		return time.Time{}, errors.Default.Wrap(err, "failed to unmarshal item")
	}
	return entity.CreatedAt, nil
}
