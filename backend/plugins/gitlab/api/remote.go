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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

type RemoteScopesChild struct {
	Type     string      `json:"type"`
	ParentId *string     `json:"parentId"`
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
}

type RemoteScopesOutput struct {
	Children      []RemoteScopesChild `json:"children"`
	NextPageToken string              `json:"nextPageToken"`
}

type SearchRemoteScopesOutput struct {
	Children []RemoteScopesChild `json:"children"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

type PageData struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Tag     string `json:"tag"`
}

const GitlabRemoteScopesPerPage int = 100
const TypeProject string = "scope"
const TypeGroup string = "group"

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input)
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
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	search, ok := input.Query["search"]
	if !ok || len(search) == 0 {
		search = []string{""}
	}

	var p int
	var err1 error
	page, ok := input.Query["page"]
	if !ok || len(page) == 0 {
		p = 1
	} else {
		p, err1 = strconv.Atoi(page[0])
		if err != nil {
			return nil, errors.BadInput.Wrap(err1, fmt.Sprintf("failed to Atoi page:%s", page[0]))
		}
	}
	var ps int
	pageSize, ok := input.Query["pageSize"]
	if !ok || len(pageSize) == 0 {
		ps = GitlabRemoteScopesPerPage
	} else {
		ps, err1 = strconv.Atoi(pageSize[0])
		if err1 != nil {
			return nil, errors.BadInput.Wrap(err1, fmt.Sprintf("failed to Atoi pageSize:%s", pageSize[0]))
		}
	}
	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	// set query
	query, err := GetQueryForSearchProject(search[0], p, ps)
	if err != nil {
		return nil, err
	}

	// request search
	res, err := apiClient.Get("search", query, nil)
	if err != nil {
		return nil, err
	}
	resBody := []tasks.GitlabApiProject{}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}

	outputBody := &SearchRemoteScopesOutput{}

	// append project to output
	for _, project := range resBody {
		child := RemoteScopesChild{
			Type:     TypeProject,
			Id:       strconv.Itoa(project.GitlabId),
			ParentId: nil,
			Name:     project.PathWithNamespace,
			Data:     tasks.ConvertProject(&project),
		}

		outputBody.Children = append(outputBody.Children, child)
	}

	outputBody.Page = p
	outputBody.PageSize = ps

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

func GetQueryFromPageData(pageData *PageData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", pageData.Page))
	query.Set("per_page", fmt.Sprintf("%v", pageData.PerPage))
	return query, nil
}

func GetQueryForSearchProject(search string, page int, perPage int) (url.Values, errors.Error) {
	query, err := GetQueryFromPageData(&PageData{Page: page, PerPage: perPage})
	if err != nil {
		return nil, err
	}
	query.Set("search", search)
	query.Set("scope", "projects")

	return query, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	projectId, _ := strconv.ParseUint(params["projectId"], 10, 64)
	return connectionId, projectId
}
