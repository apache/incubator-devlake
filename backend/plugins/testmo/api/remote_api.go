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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
	"github.com/apache/incubator-devlake/plugins/testmo/models/raw"
)

type TestmoRemotePagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type ProjectPage struct {
	Result   []raw.TestmoProject `json:"result"`
	Page     *int                `json:"page"`
	PrevPage *int                `json:"prev_page"`
	NextPage *int                `json:"next_page"`
	LastPage *int                `json:"last_page"`
	PerPage  *int                `json:"per_page"`
	Total    *int                `json:"total"`
}

func queryTestmoRemoteScopes(
	apiClient plugin.ApiClient,
	_ string,
	page TestmoRemotePagination,
	search string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TestmoProject],
	nextPage *TestmoRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	var res *http.Response
	queryParams := url.Values{
		"page": {strconv.Itoa(page.Page)},
	}
	if search != "" {
		queryParams["search"] = []string{search}
	}

	res, err = apiClient.Get("projects", queryParams, nil)
	if err != nil {
		return
	}

	// Parse as ProjectPage response
	response := &ProjectPage{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return
	}
	// append project to output
	for _, project := range response.Result {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.TestmoProject]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       strconv.FormatUint(project.Id, 10),
			Name:     project.Name,
			FullName: project.Name,
			Data: &models.TestmoProject{
				Id:                           project.Id,
				Name:                         project.Name,
				IsCompleted:                  project.IsCompleted,
				MilestoneCount:               project.MilestoneCount,
				MilestoneActiveCount:         project.MilestoneActiveCount,
				MilestoneCompletedCount:      project.MilestoneCompletedCount,
				RunCount:                     project.RunCount,
				RunActiveCount:               project.RunActiveCount,
				RunClosedCount:               project.RunClosedCount,
				AutomationSourceCount:        project.AutomationSourceCount,
				AutomationSourceActiveCount:  project.AutomationSourceActiveCount,
				AutomationSourceRetiredCount: project.AutomationSourceRetiredCount,
				AutomationRunCount:           project.AutomationRunCount,
				AutomationRunActiveCount:     project.AutomationRunActiveCount,
				AutomationRunCompletedCount:  project.AutomationRunCompletedCount,
			},
		})
	}

	// Check if there are more pages
	if response.NextPage != nil {
		perPage := 100 // Default per_page value used by Testmo API
		if response.PerPage != nil {
			perPage = *response.PerPage
		}
		nextPage = &TestmoRemotePagination{
			Page:    *response.NextPage,
			PerPage: perPage,
		}
	}

	return
}

func listTestmoRemoteScopes(
	connection *models.TestmoConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page TestmoRemotePagination,
) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.TestmoProject],
	*TestmoRemotePagination,
	errors.Error,
) {
	return queryTestmoRemoteScopes(apiClient, groupId, page, "")
}

func searchTestmoRemoteScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TestmoProject],
	err errors.Error,
) {
	children, _, err = queryTestmoRemoteScopes(apiClient, "", TestmoRemotePagination{
		Page:    params.Page,
		PerPage: params.PageSize,
	}, params.Search)
	return
}

// RemoteScopes list all available scopes (projects) for this connection
// @Summary list all available scopes (projects) for this connection
// @Description list all available scopes (projects) for this connection
// @Tags plugins/testmo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/testmo/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/testmo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/testmo/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}
