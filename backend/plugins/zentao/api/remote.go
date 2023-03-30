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
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

type ProductResponse struct {
	Limit  int                       `json:"limit"`
	Page   int                       `json:"page"`
	Total  int                       `json:"total"`
	Values []models.ZentaoProductRes `json:"products"`
}

type ProjectResponse struct {
	Limit  int                    `json:"limit"`
	Page   int                    `json:"page"`
	Total  int                    `json:"total"`
	Values []models.ZentaoProject `json:"projects"`
}

func getGroup(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.ZentaoConnection) ([]api.BaseRemoteGroupResponse, errors.Error) {
	return []api.BaseRemoteGroupResponse{
		{
			Id:   `products`,
			Name: `Products`,
		},
		{
			Id:   `projects`,
			Name: `Projects`,
		},
	}, nil
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	groupId, ok := input.Query["groupId"]
	if !ok || len(groupId) == 0 {
		groupId = []string{""}
	}
	gid := groupId[0]
	if gid == "" {
		return productRemoteHelper.GetScopesFromRemote(input, getGroup, nil)
	} else if gid == `products` {
		return productRemoteHelper.GetScopesFromRemote(input,
			nil,
			func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.ZentaoConnection) ([]models.ZentaoProductRes, errors.Error) {
				query := initialQuery(queryData)
				// create api client
				apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
				if err != nil {
					return nil, err
				}

				query.Set("sort", "name")
				// list projects part
				res, err := apiClient.Get("/products", query, nil)
				if err != nil {
					return nil, err
				}

				resBody := &ProductResponse{}
				err = api.UnmarshalResponse(res, resBody)
				if err != nil {
					return nil, err
				}
				return resBody.Values, nil
			})
	} else if gid == `projects` {
		return projectRemoteHelper.GetScopesFromRemote(input,
			nil,
			func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.ZentaoConnection) ([]models.ZentaoProject, errors.Error) {
				query := initialQuery(queryData)
				// create api client
				apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
				if err != nil {
					return nil, err
				}

				query.Set("sort", "name")
				// list projects part
				res, err := apiClient.Get("/projects", query, nil)
				if err != nil {
					return nil, err
				}

				resBody := &ProjectResponse{}
				err = api.UnmarshalResponse(res, resBody)
				if err != nil {
					return nil, err
				}
				return resBody.Values, nil
			})
	}
	return nil, nil
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", queryData.Page))
	query.Set("limit", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
