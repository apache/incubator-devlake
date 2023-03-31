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
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"net/http"
	"net/url"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		nil,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BambooConnection) ([]models.ApiBambooProject, errors.Error) {
			query := initialQuery(queryData)
			// create api client
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, err
			}
			res, err := apiClient.Get("/project.json", query, nil)

			if err != nil {
				return nil, err
			}

			resBody := models.ApiBambooProjectResponse{}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody.Projects.Projects, err
		})
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context2.BasicRes, queryData *api.RemoteQueryData, connection models.BambooConnection) ([]models.ApiBambooProject, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			query.Set("searchTerm", queryData.Search[0])
			// request search
			res, err := apiClient.Get("search/projects.json", query, nil)
			if err != nil {
				return nil, err
			}
			resBody := models.ApiBambooSearchProjectResponse{}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			var apiBambooProjects []models.ApiBambooProject
			// append project to output
			for _, apiResult := range resBody.SearchResults {
				apiProject, err := GetApiProject(apiResult.SearchEntity.Key, apiClient)
				if err != nil {
					return nil, err
				}

				apiBambooProjects = append(apiBambooProjects, *apiProject)
			}
			return apiBambooProjects, err
		})
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("showEmpty", fmt.Sprintf("%v", true))
	query.Set("max-result", fmt.Sprintf("%v", queryData.PerPage))
	query.Set("start-index", fmt.Sprintf("%v", (queryData.Page-1)*queryData.PerPage))
	return query
}

// move from blueprint_v200 because of cycle import
func GetApiProject(
	projectKey string,
	apiClient aha.ApiClientAbstract,
) (*models.ApiBambooProject, errors.Error) {
	projectRes := &models.ApiBambooProject{}
	res, err := apiClient.Get(fmt.Sprintf("project/%s.json", projectKey), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting project detail from %s", res.Request.URL.String()))
	}
	err = api.UnmarshalResponse(res, projectRes)
	if err != nil {
		return nil, err
	}
	return projectRes, nil
}
