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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/models/common"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models/raw"
)

type PagerdutyRemotePagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ServiceResponse struct {
	Offset   int           `json:"offset"`
	Limit    int           `json:"limit"`
	More     bool          `json:"more"`
	Total    int           `json:"total"`
	Services []raw.Service `json:"services"`
}

func queryPagerdutyRemoteScopes(
	apiClient plugin.ApiClient,
	_ string,
	page PagerdutyRemotePagination,
	search string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.Service],
	nextPage *PagerdutyRemotePagination,
	err errors.Error,
) {
	if page.Limit == 0 {
		page.Limit = 50
	}
	var res *http.Response
	res, err = apiClient.Get("/services", url.Values{
		"offset": {strconv.Itoa(page.Offset)},
		"limit":  {strconv.Itoa(page.Limit)},
		"query":  {search},
	}, nil)
	if err != nil {
		return
	}
	response := &ServiceResponse{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return
	}
	// append service to output
	for _, service := range response.Services {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.Service]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       service.Id,
			Name:     service.Name,
			FullName: service.Name,
			Data: &models.Service{
				Url:  service.HtmlUrl,
				Id:   service.Id,
				Name: service.Name,
				Scope: common.Scope{
					NoPKModel: common.NoPKModel{
						CreatedAt: service.CreatedAt,
					},
				},
			},
		})
	}

	// check service count
	if response.More {
		nextPage = &PagerdutyRemotePagination{
			Offset: page.Offset + page.Limit,
			Limit:  page.Limit,
		}
	}

	return
}

func listPagerdutyRemoteScopes(
	connection *models.PagerDutyConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page PagerdutyRemotePagination,
) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.Service],
	*PagerdutyRemotePagination,
	errors.Error,
) {
	return queryPagerdutyRemoteScopes(apiClient, groupId, page, "")
}

func searchPagerdutyRemoteScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.Service],
	err errors.Error,
) {
	children, _, err = queryPagerdutyRemoteScopes(apiClient, "", PagerdutyRemotePagination{
		Offset: (params.Page - 1) * params.PageSize,
		Limit:  params.PageSize,
	}, params.Search)
	return
}

// RemoteScopes list all available scopes (services) for this connection
// @Summary list all available scopes (services) for this connection
// @Description list all available scopes (services) for this connection
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/github
// @Router /plugins/github/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
