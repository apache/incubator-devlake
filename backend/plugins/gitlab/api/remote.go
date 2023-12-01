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
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

type GitlabRemotePagination struct {
	Page    int    `json:"page" mapstructure:"page"`
	PerPage int    `json:"per_page" mapstructure:"per_page"`
	Step    string `json:"step" mapstructure:"step"`
}

func (p GitlabRemotePagination) ToQuery() url.Values {
	return url.Values{
		"page":     {fmt.Sprintf("%v", p.Page)},
		"per_page": {fmt.Sprintf("%v", p.PerPage)},
	}
}

func listGitlabRemoteScopes(
	connection *models.GitlabConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page GitlabRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject],
	nextPage *GitlabRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	if page.PerPage == 0 {
		page.PerPage = 100
	}
	if page.Step == "" {
		page.Step = "group"
	}

	// load all groups unless groupId is user's own account
	if page.Step == "group" && !strings.HasPrefix(groupId, "users/") {
		children, nextPage, err = listGitlabRemoteGroups(connection, apiClient, groupId, page)
		if err != nil {
			return
		}
		// no more groups
		if nextPage == nil {
			nextPage = &GitlabRemotePagination{
				Page:    1,
				PerPage: page.PerPage,
				Step:    "project",
			}
		}
	} else {
		// load all project under the group or user's own account
		children, nextPage, err = listGitlabRemoteProjects(connection, apiClient, groupId, page)
	}
	return
}

func listGitlabRemoteGroups(
	connection *models.GitlabConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page GitlabRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject],
	nextPage *GitlabRemotePagination,
	err errors.Error,
) {
	apiPath := ""
	query := page.ToQuery()
	var res *http.Response
	if groupId == "" && page.Page == 1 {
		// make users own account as a group
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       fmt.Sprintf("users/%v", apiClient.GetData("UserId")),
			Name:     apiClient.GetData("UserName").(string),
			FullName: apiClient.GetData("UserName").(string),
		})
	}
	if groupId == "" {
		apiPath = "groups"
		query.Set("top_level_only", "true")
	} else {
		apiPath = fmt.Sprintf("groups/%s/subgroups", groupId)
	}
	res, err = apiClient.Get(apiPath, query, nil)
	var resGroups []models.GroupResponse
	errors.Must(api.UnmarshalResponse(res, &resGroups))
	for _, group := range resGroups {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       fmt.Sprintf("%v", group.Id),
			Name:     group.Name,
			FullName: group.FullPath,
		})
	}
	nextPage = getNextPage(&page, res)
	return
}

func listGitlabRemoteProjects(
	connection *models.GitlabConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page GitlabRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject],
	nextPage *GitlabRemotePagination,
	err errors.Error,
) {
	apiPath := ""
	query := page.ToQuery()
	query.Set("archived", "false")
	query.Set("min_access_level", "20")
	//
	if strings.HasPrefix(groupId, "users/") {
		apiPath = fmt.Sprintf("%s/projects", groupId)
	} else {
		apiPath = fmt.Sprintf("/groups/%s/projects", groupId)
	}
	res, err := apiClient.Get(apiPath, query, nil)
	var resProjects []models.GitlabApiProject
	errors.Must(api.UnmarshalResponse(res, &resProjects))
	for _, project := range resProjects {
		children = append(children, toProjectModel(&project))
	}
	nextPage = getNextPage(&page, res)
	return
}

func toProjectModel(project *models.GitlabApiProject) dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject] {
	return dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]{
		Type:     api.RAS_ENTRY_TYPE_SCOPE,
		Id:       fmt.Sprintf("%v", project.GitlabId),
		Name:     project.Name,
		FullName: project.PathWithNamespace,
	}
}

func getNextPage(page *GitlabRemotePagination, res *http.Response) *GitlabRemotePagination {
	if res.Header.Get("x-next-page") == "" {
		return nil
	}
	return &GitlabRemotePagination{
		Page:    page.Page + 1,
		PerPage: page.PerPage,
		Step:    page.Step,
	}
}

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
	return raScopeList.Get(input)
}

func searchGitlabScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject],
	err errors.Error,
) {
	res, err := apiClient.Get(
		"projects",
		url.Values{
			"search":           []string{params.Search},
			"page":             []string{fmt.Sprintf("%v", params.Page)},
			"per_page":         []string{fmt.Sprintf("%v", params.PageSize)},
			"archived":         {"false"},
			"min_access_level": {"20"},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var resBody []models.GitlabApiProject
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	for _, r := range resBody {
		children = append(children, toProjectModel(&r))
	}
	return
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
	return raScopeSearch.Get(input)
}

// Proxy is a proxy to Gitlab API
// @Summary Proxy to Gitlab API
// @Description Proxy to Gitlab API
// @Tags plugins/gitlab
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to Gitlab API"
// @Router /plugins/gitlab/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
