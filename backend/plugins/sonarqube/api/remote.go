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
	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"net/url"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		nil,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.SonarqubeConnection) ([]models.SonarqubeApiProject, errors.Error) {
			query := initialQuery(queryData)
			// create api client
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, err
			}

			res, err := apiClient.Get("projects/search", query, nil)
			if err != nil {
				return nil, err
			}

			var resBody struct {
				Data []models.SonarqubeApiProject `json:"components"`
			}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			return resBody.Data, nil
		})
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search keyword"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context2.BasicRes, queryData *api.RemoteQueryData, connection models.SonarqubeConnection) ([]models.SonarqubeApiProject, errors.Error) {
			query := initialQuery(queryData)
			query.Set("q", queryData.Search[0])
			// create api client
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, err
			}

			// request search
			res, err := apiClient.Get("projects/search", query, nil)
			if err != nil {
				return nil, err
			}
			var resBody struct {
				Data []models.SonarqubeApiProject `json:"components"`
			}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			return resBody.Data, nil
		})
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("p", fmt.Sprintf("%v", queryData.Page))
	query.Set("ps", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
