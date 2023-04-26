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
	"net/http"
	"net/url"

	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.GitlabConnection) ([]models.GroupResponse, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			var res *http.Response
			if gid == "" {
				query.Set("top_level_only", "true")
				res, err = apiClient.Get("groups", query, nil)
				if err != nil {
					return nil, err
				}
			} else {
				if gid[:6] == "group:" {
					gid = gid[6:]
				}
				res, err = apiClient.Get(fmt.Sprintf("groups/%s/subgroups", gid), query, nil)
				if err != nil {
					return nil, err
				}
			}
			var resBody []models.GroupResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			return resBody, err
		},
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.GitlabConnection) ([]models.GitlabApiProject, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			var res *http.Response
			if gid == "" {
				res, err = apiClient.Get(fmt.Sprintf("users/%d/projects", apiClient.GetData("UserId")), query, nil)
				if err != nil {
					return nil, err
				}
			} else {
				query.Set("with_shared", "false")
				if gid[:6] == "group:" {
					gid = gid[6:]
				}
				res, err = apiClient.Get(fmt.Sprintf("/groups/%s/projects", gid), query, nil)
				if err != nil {
					return nil, err
				}
			}
			var resBody []models.GitlabApiProject
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody, err
		})
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context2.BasicRes, queryData *api.RemoteQueryData, connection models.GitlabConnection) ([]models.GitlabApiProject, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			query.Set("search", queryData.Search[0])
			query.Set("scope", "projects")
			// request search
			res, err := apiClient.Get("search", query, nil)
			if err != nil {
				return nil, err
			}
			var resBody []models.GitlabApiProject
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			for i := 0; i < len(resBody); i++ {
				// as we need to set PathWithNamespace to name in SearchRemoteScopes, but interface.ScopeName will return name, so we switch it
				resBody[i].Name, resBody[i].PathWithNamespace = resBody[i].PathWithNamespace, resBody[i].Name
			}
			return resBody, err
		})
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", queryData.Page))
	query.Set("per_page", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
