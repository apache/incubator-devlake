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

const USERS_PREFIX = "user:"

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
	if page.Step == "group" && !strings.HasPrefix(groupId, USERS_PREFIX) {
		children, nextPage, err = listGitlabRemoteGroups(connection, apiClient, groupId, page)
		if err != nil {
			return
		}
	}
	if groupId == "" || nextPage != nil {
		return
	}
	// no more groups, start to load projects under the group
	var moreChild []dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]
	moreChild, nextPage, err = listGitlabRemoteProjects(connection, apiClient, groupId, GitlabRemotePagination{
		Page:    page.Page,
		PerPage: page.PerPage,
		Step:    "project",
	})
	if err != nil {
		return
	}
	children = append(children, moreChild...)
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
			Id:       USERS_PREFIX + fmt.Sprintf("%v", apiClient.GetData("UserId")),
			Name:     apiClient.GetData("UserName").(string),
			FullName: apiClient.GetData("UserName").(string),
		})
	}
	var parentId *string
	if groupId == "" {
		apiPath = "groups"
		query.Set("top_level_only", "true")
	} else {
		apiPath = fmt.Sprintf("groups/%s/subgroups", groupId)
		parentId = &groupId
	}
	res, err = apiClient.Get(apiPath, query, nil)
	var resGroups []models.GroupResponse
	errors.Must(api.UnmarshalResponse(res, &resGroups))
	for _, group := range resGroups {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       fmt.Sprintf("%v", group.Id),
			ParentId: parentId,
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
	if strings.HasPrefix(groupId, USERS_PREFIX) {
		apiPath = fmt.Sprintf("users/%s/projects", strings.TrimPrefix(groupId, USERS_PREFIX))
	} else {
		apiPath = fmt.Sprintf("groups/%s/projects", groupId)
	}
	res, err := apiClient.Get(apiPath, query, nil)
	if err != nil {
		return nil, nil, err
	}
	var resProjects []models.GitlabApiProject
	errors.Must(api.UnmarshalResponse(res, &resProjects))
	for _, project := range resProjects {
		children = append(children, toProjectModel(&project))
	}
	nextPage = getNextPage(&page, res)
	return
}

func toProjectModel(project *models.GitlabApiProject) dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject] {
	var parentId string
	if project.Namespace.Kind == "user" {
		parentId = USERS_PREFIX + fmt.Sprintf("%v", project.Owner.ID)
	} else {
		parentId = fmt.Sprintf("%v", project.Namespace.ID)
	}
	return dsmodels.DsRemoteApiScopeListEntry[models.GitlabProject]{
		Type:     api.RAS_ENTRY_TYPE_SCOPE,
		Id:       fmt.Sprintf("%v", project.GitlabId),
		ParentId: &parentId,
		Name:     project.Name,
		FullName: project.PathWithNamespace,
		Data:     project.ConvertApiScope(),
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

// RemoteScopes list all available scopes on the remote server
// @Summary list all available scopes on the remote server
// @Description list all available scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GitlabProject]
// @Tags plugins/gitlab
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

// SearchRemoteScopes searches projects on the remote server
// @Summary searches projects on the remote server
// @Description searches projects on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GitlabProject] "the parentIds are always null"
// @Tags plugins/gitlab
// @Router /plugins/gitlab/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Router /plugins/gitlab/connections/{connectionId}/proxy/{path} [GET]
// @Tags plugins/gitlab
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
