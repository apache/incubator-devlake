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
	"io"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const rawCopilotSeatsTable = "copilot_seats"

func parseCopilotSeatsFromResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	if res == nil {
		return nil, errors.Default.New("res is nil")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error reading response body of %s", res.Request.URL.String()))
	}

	// Some GitHub endpoints return a top-level array, others return an object wrapper.
	// Support both, returning the underlying array of seat objects.
	var rawMessages []json.RawMessage
	if jsonErr := json.Unmarshal(body, &rawMessages); jsonErr == nil {
		return rawMessages, nil
	}

	var wrapped struct {
		Seats []json.RawMessage `json:"seats"`
	}
	if jsonErr := json.Unmarshal(body, &wrapped); jsonErr != nil {
		return nil, errors.Default.Wrap(errors.Convert(jsonErr), fmt.Sprintf("error decoding response of %s: raw response: %s", res.Request.URL.String(), string(body)))
	}
	return wrapped.Seats, nil
}

func CollectCopilotSeatAssignments(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not CopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	apiClient, err := CreateApiClient(taskCtx.TaskContext(), connection)
	if err != nil {
		return err
	}

	rawArgs := helper.RawDataSubTaskArgs{
		Ctx:   taskCtx,
		Table: rawCopilotSeatsTable,
		Options: copilotRawParams{
			ConnectionId: data.Options.ConnectionId,
			ScopeId:      data.Options.ScopeId,
			Organization: connection.Organization,
			Endpoint:     connection.Endpoint,
		},
	}

	collector, err := helper.NewStatefulApiCollector(rawArgs)
	if err != nil {
		return err
	}

	perPage := 100
	err = collector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   apiClient,
		PageSize:    perPage,
		UrlTemplate: fmt.Sprintf("orgs/%s/copilot/billing/seats", connection.Organization),
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			q := url.Values{}
			q.Set("per_page", fmt.Sprintf("%d", reqData.Pager.Size))
			q.Set("page", fmt.Sprintf("%d", reqData.Pager.Page))
			return q, nil
		},
		GetNextPageCustomData: func(prevReqData *helper.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			// Standard page/per_page pagination; nothing extra to carry between pages.
			return nil, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			if res.StatusCode >= 400 {
				body, _ := io.ReadAll(res.Body)
				res.Body.Close()
				return nil, buildGitHubApiError(res.StatusCode, connection.Organization, body, res.Header.Get("Retry-After"))
			}
			return parseCopilotSeatsFromResponse(res)
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}
