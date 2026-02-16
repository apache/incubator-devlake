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
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// GhCopilotRemotePagination is a placeholder for scope list pagination.
// Copilot scopes are organization-level and currently return a single entry.
type GhCopilotRemotePagination struct {
	Page int `json:"page"`
}

func listGhCopilotRemoteScopes(
	connection *models.GhCopilotConnection,
	_ plugin.ApiClient,
	_ string,
	_ GhCopilotRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope],
	nextPage *GhCopilotRemotePagination,
	err errors.Error,
) {
	if connection == nil {
		return nil, nil, errors.BadInput.New("connection is required")
	}
	organization := strings.TrimSpace(connection.Organization)

	if connection.HasEnterprise() {
		enterprise := strings.TrimSpace(connection.Enterprise)
		if enterprise != "" {
			scopeId := enterprise
			if organization != "" {
				scopeId = enterprise + "/" + organization
			}
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope]{
				Type:     api.RAS_ENTRY_TYPE_SCOPE,
				Id:       scopeId,
				Name:     scopeId,
				FullName: scopeId,
				Data: &models.GhCopilotScope{
					Id:           scopeId,
					Organization: organization,
					Enterprise:   enterprise,
					Name:         scopeId,
					FullName:     scopeId,
				},
			})
		}
	} else if organization != "" {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       organization,
			Name:     organization,
			FullName: organization,
			Data: &models.GhCopilotScope{
				Id:           organization,
				Organization: organization,
				Name:         organization,
				FullName:     organization,
			},
		})
	}

	return children, nil, nil
}

func searchGhCopilotRemoteScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope],
	err errors.Error,
) {
	if params == nil {
		return []dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope]{}, nil
	}
	query := strings.TrimSpace(params.Search)
	if query == "" {
		return []dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope]{}, nil
	}
	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	resp, err := apiClient.Get(
		"user/orgs",
		url.Values{
			"page":     []string{strconv.Itoa(page)},
			"per_page": []string{strconv.Itoa(pageSize)},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	var orgs []struct {
		Login string `json:"login"`
	}
	if err := api.UnmarshalResponse(resp, &orgs); err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	for _, org := range orgs {
		orgName := strings.TrimSpace(org.Login)
		if orgName == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(orgName), queryLower) {
			continue
		}
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GhCopilotScope]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       orgName,
			Name:     orgName,
			FullName: orgName,
			Data: &models.GhCopilotScope{
				Id:           orgName,
				Organization: orgName,
				Name:         orgName,
				FullName:     orgName,
			},
		})
	}

	return children, nil
}

// RemoteScopes list all available scopes (organizations) for this connection
// @Summary list all available scopes (organizations) for this connection
// @Description list all available scopes (organizations) for this connection
// @Tags plugins/gh-copilot
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GhCopilotScope]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gh-copilot/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes searches organization scopes for this connection
// @Summary searches organization scopes for this connection
// @Description searches organization scopes for this connection
// @Tags plugins/gh-copilot
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GhCopilotScope] "the parentIds are always null"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gh-copilot/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}
