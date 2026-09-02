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
	"strings"
	"time"

	"github.com/apache/devlake/core/models/common"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/devlake/helpers/pluginhelper/api/models"
	"github.com/apache/devlake/plugins/incidentio/models"
)

type IncidentioRemotePagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type IncidentTypesResponse struct {
	IncidentTypes []struct {
		Id          string     `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`
		CreatedAt   *time.Time `json:"created_at"`
	} `json:"incident_types"`
}

// queryIncidentioRemoteScopes lists incident types as scopes. The
// endpoint is not paginated and returns the full list, so search is
// applied client-side and there is never a next page.
func queryIncidentioRemoteScopes(
	apiClient plugin.ApiClient,
	_ string,
	page IncidentioRemotePagination,
	search string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.IncidentType],
	nextPage *IncidentioRemotePagination,
	err errors.Error,
) {
	var res *http.Response
	res, err = apiClient.Get("v1/incident_types", nil, nil)
	if err != nil {
		return
	}
	response := &IncidentTypesResponse{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return
	}
	for _, item := range response.IncidentTypes {
		if search != "" && !strings.Contains(strings.ToLower(item.Name), strings.ToLower(search)) {
			continue
		}
		entry := dsmodels.DsRemoteApiScopeListEntry[models.IncidentType]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       item.Id,
			Name:     item.Name,
			FullName: item.Name,
			Data: &models.IncidentType{
				Id:   item.Id,
				Name: item.Name,
				Scope: common.Scope{
					NoPKModel: common.NoPKModel{},
				},
			},
		}
		if item.CreatedAt != nil {
			entry.Data.Scope.NoPKModel.CreatedAt = *item.CreatedAt
		}
		children = append(children, entry)
	}

	return
}

func listIncidentioRemoteScopes(
	connection *models.IncidentioConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page IncidentioRemotePagination,
) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.IncidentType],
	*IncidentioRemotePagination,
	errors.Error,
) {
	return queryIncidentioRemoteScopes(apiClient, groupId, page, "")
}

func searchIncidentioRemoteScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.IncidentType],
	err errors.Error,
) {
	children, _, err = queryIncidentioRemoteScopes(apiClient, "", IncidentioRemotePagination{
		Page:    params.Page,
		PerPage: params.PageSize,
	}, params.Search)
	return
}

// RemoteScopes list all available scopes (incident types) for this connection
// @Summary list all available scopes (incident types) for this connection
// @Description list all available scopes (incident types) for this connection
// @Tags plugins/incidentio
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/incidentio/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/incidentio
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/incidentio/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}
