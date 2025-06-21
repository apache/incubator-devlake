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
	"fmt"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
)

type TeambitionPagination struct {
	PageToken string `json:"pageToken"`
	PageSize  int    `json:"pageSize"`
}

func queryTeambitionProjects(
	apiClient plugin.ApiClient,
	keyword string,
	page TeambitionPagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TeambitionProject],
	nextPage *TeambitionPagination,
	err errors.Error,
) {
	if page.PageSize == 0 {
		page.PageSize = 50
	}
	res, err := apiClient.Get("v3/project/query", url.Values{
		"name":      {keyword},
		"pageSize":  {fmt.Sprintf("%v", page.PageSize)},
		"pageToken": {page.PageToken},
	}, nil)
	if err != nil {
		return
	}
	resBody := struct {
		Result        []models.TeambitionProject `json:"result"`
		NextPageToken string                     `json:"nextPageToken"`
	}{}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return
	}
	for _, project := range resBody.Result {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.TeambitionProject]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       fmt.Sprintf("%v", project.Id),
			ParentId: nil,
			Name:     project.Name,
			FullName: project.Name,
			Data:     &project,
		})
	}
	if resBody.NextPageToken != "" {
		nextPage = &TeambitionPagination{
			PageToken: resBody.NextPageToken,
			PageSize:  page.PageSize,
		}
	}
	return
}

func listTeambitionRemoteScopes(
	connection *models.TeambitionConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page TeambitionPagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TeambitionProject],
	nextPage *TeambitionPagination,
	err errors.Error,
) {
	// construct the query and request
	return queryTeambitionProjects(apiClient, "", page)
}

func searchTeambitionRemoteProjects(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TeambitionProject],
	err errors.Error,
) {
	if params.Page == 0 {
		params.Page = 1
	}
	page := TeambitionPagination{
		PageSize: params.PageSize,
	}
	children, _, err = queryTeambitionProjects(apiClient, params.Search, page)
	return
}

type Entry = dsmodels.DsRemoteApiScopeListEntry[models.TeambitionProject]
type Node struct {
	entry *Entry
}
type Children []*Node

func (a Children) Len() int { return len(a) }
func (a Children) Less(i, j int) bool {
	if a[i].entry.Type != a[j].entry.Type {
		return a[i].entry.Type < a[j].entry.Type
	}
	return a[i].entry.Name < a[j].entry.Name
}
func (a Children) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.TapdWorkspace]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes searches scopes on the remote server
// @Summary searches scopes on the remote server
// @Description searches scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.SonarqubeProject] "the parentIds are always null"
// @Tags plugins/sonarqube
// @Router /plugins/sonarqube/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/github
// @Router /plugins/teambition/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
