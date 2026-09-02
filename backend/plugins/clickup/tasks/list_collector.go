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
	"strconv"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
)

// RAW_LIST_TABLE holds the raw lists of a scoped folder (GET /folder/{id}/list).
const RAW_LIST_TABLE = "clickup_lists"

// clickUpListResponse mirrors GET /folder/{id}/list. The endpoint is not
// paginated and returns the folder's lists under `lists`.
type clickUpListResponse struct {
	Lists []json.RawMessage `json:"lists"`
}

// archivedInput drives the collector to request the folder's lists twice: once
// for active lists and once for archived lists (ClickUp's /folder/{id}/list
// returns active by default and archived only when archived=true).
type archivedInput struct {
	Archived bool
}

// archivedIterator yields {false, true} so a single Execute collects both the
// active and the archived lists into the same raw table.
type archivedIterator struct {
	vals []bool
	i    int
}

func newArchivedIterator() *archivedIterator { return &archivedIterator{vals: []bool{false, true}} }

func (it *archivedIterator) HasNext() bool { return it.i < len(it.vals) }

func (it *archivedIterator) Fetch() (interface{}, errors.Error) {
	v := it.vals[it.i]
	it.i++
	return &archivedInput{Archived: v}, nil
}

func (it *archivedIterator) Close() errors.Error { return nil }

var CollectListMeta = plugin.SubTaskMeta{
	Name:             "Collect Lists",
	EntryPoint:       CollectLists,
	EnabledByDefault: true,
	Description:      "Collect the lists (backlog + sprint lists, active and archived) of a scoped ClickUp folder",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

var _ plugin.SubTaskEntryPoint = CollectLists

func CollectLists(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_LIST_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       newArchivedIterator(),
		UrlTemplate: "folder/{{ .Params.FolderId }}/list",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			archived := false
			if in, ok := reqData.Input.(*archivedInput); ok {
				archived = in.Archived
			}
			query.Set("archived", strconv.FormatBool(archived))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resp clickUpListResponse
			if err := api.UnmarshalResponse(res, &resp); err != nil {
				return nil, err
			}
			return resp.Lists, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}
