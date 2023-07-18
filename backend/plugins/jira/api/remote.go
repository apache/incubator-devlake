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

package api

import (
	"context"
	"fmt"
	"net/url"

	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		nil,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.JiraConnection) ([]apiv2models.Board, errors.Error) {
			query := initialQuery(queryData)
			// create api client
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, err
			}
			res, err := apiClient.Get("agile/1.0/board", query, nil)
			if err != nil {
				return nil, err
			}

			resBody := struct {
				MaxResults int                 `json:"maxResults"`
				StartAt    int                 `json:"startAt"`
				Values     []apiv2models.Board `json:"values"`
			}{}

			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			if (queryData.PerPage != resBody.MaxResults) ||
				(((queryData.Page - 1) * queryData.PerPage) != resBody.StartAt) {
				analyzingQuery(resBody.MaxResults, resBody.StartAt, queryData)
			}

			return resBody.Values, err
		})
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("maxResults", fmt.Sprintf("%v", queryData.PerPage))
	query.Set("startAt", fmt.Sprintf("%v", (queryData.Page-1)*queryData.PerPage))
	return query
}

func analyzingQuery(maxResults int, startAt int, queryData *api.RemoteQueryData) {
	if maxResults != 0 {
		queryData.PerPage = maxResults
		queryData.Page = startAt/maxResults + 1
	}
}
