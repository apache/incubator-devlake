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
	gocontext "context"
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
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
	return remoteHelper.GetScopesFromRemote(input, nil, getRemotePlans)
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
		func(basicRes context.BasicRes, queryData *api.RemoteQueryData, connection models.BambooConnection) ([]models.ApiBambooPlan, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			if len(queryData.Search) == 0 {
				return nil, errors.BadInput.New("empty search query")
			}
			query.Set("searchTerm", queryData.Search[0])
			// request search
			res, err := apiClient.Get("search/plans.json", query, nil)
			if err != nil {
				return nil, err
			}

			resBody := models.ApiBambooSearchPlanResponse{}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			var apiBambooPlans []models.ApiBambooPlan
			// append project to output
			for _, apiResult := range resBody.SearchResults {
				bambooPlan := models.ApiBambooPlan{
					Key:  apiResult.SearchEntity.Key,
					Name: apiResult.SearchEntity.Name(),
				}
				apiBambooPlans = append(apiBambooPlans, bambooPlan)
			}
			return apiBambooPlans, err
		})
}

func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.BambooConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return remoteHelper.ProxyApiGet(connection, input.Params["path"], input.Query)
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("showEmpty", fmt.Sprintf("%v", false))
	query.Set("max-result", fmt.Sprintf("%v", queryData.PerPage))
	query.Set("start-index", fmt.Sprintf("%v", (queryData.Page-1)*queryData.PerPage))
	return query
}

func getRemotePlans(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BambooConnection) ([]models.ApiBambooPlan, errors.Error) {
	query := initialQuery(queryData)
	// create api client
	apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("plan.json", query, nil)

	if err != nil {
		return nil, err
	}
	var planRes struct {
		Expand string `json:"expand"`
		Link   struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"link"`
		Plans struct {
			Size       int                    `json:"size"`
			Expand     string                 `json:"expand"`
			StartIndex int                    `json:"start-index"`
			MaxResult  int                    `json:"max-result"`
			Plan       []models.ApiBambooPlan `json:"plan"`
		} `json:"plans"`
	}
	err = api.UnmarshalResponse(res, &planRes)
	if err != nil {
		return nil, err
	}
	return planRes.Plans.Plan, err
}
