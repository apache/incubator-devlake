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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_STORY_TABLE = "zentao_api_stories"

var _ plugin.SubTaskEntryPoint = CollectStory

type storyInput struct {
	Path        string
	ProjectId   int64
	ProductId   int64
	ExecutionId int64
}

func CollectStory(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	// project iterator
	iter0 := newIteratorFromSlice([]interface{}{&storyInput{ProjectId: data.Options.ProjectId, Path: fmt.Sprintf("/projects/%d", data.Options.ProjectId)}})

	// product iterator
	productCursor, productIterator, err := getProductIterator(taskCtx)
	if err != nil {
		return err
	}
	defer productCursor.Close()
	iter1 := newIteratorWrapper(productIterator, func(arg interface{}) interface{} {
		return &storyInput{ProductId: arg.(*input).Id, Path: fmt.Sprintf("/products/%d", arg.(*input).Id)}
	})

	// execution iterator
	executionCursor, executionIterator, err := getExecutionIterator(taskCtx)
	if err != nil {
		return err
	}
	defer executionCursor.Close()
	iter2 := newIteratorWrapper(executionIterator, func(arg interface{}) interface{} {
		return &storyInput{ExecutionId: arg.(*input).Id, Path: fmt.Sprintf("/executions/%d", arg.(*input).Id)}
	})

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_STORY_TABLE,
		},
		Input:       newIteratorConcator(iter0, iter1, iter2),
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "{{ .Input.Path }}/stories",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("status", "allstory")
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Story []json.RawMessage `json:"stories"`
			}
			err := api.UnmarshalResponse(res, &data)
			if errors.Is(err, api.ErrEmptyResponse) {
				return nil, nil
			}
			if err != nil {
				return nil, errors.Default.Wrap(err, "error reading endpoint response by Zentao bug collector")
			}
			return data.Story, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectStoryMeta = plugin.SubTaskMeta{
	Name:             "collectStory",
	EntryPoint:       CollectStory,
	EnabledByDefault: true,
	Description:      "Collect Story data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
