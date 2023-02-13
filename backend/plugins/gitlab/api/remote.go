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
	"encoding/base64"
	"encoding/json"
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
	Data     interface{} `json:"data"`
}

type RemoteScopesOutput struct {
	Children      []RemoteScopesChild `json:"children"`
	NextPageToken string              `json:"nextPageToken"`
}

type PageData struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Tag     string `json:"tag"`
}

type GroupResponse struct {
	Id                   int    `json:"id"`
	WebUrl               string `json:"web_url"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	Description          string `json:"description"`
	Visibility           string `json:"visibility"`
	LfsEnabled           bool   `json:"lfs_enabled"`
	AvatarUrl            string `json:"avatar_url"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
	FullName             string `json:"full_name"`
	FullPath             string `json:"full_path"`
	ParentId             *int   `json:"parent_id"`
	LdapCN               string `json:"ldap_cn"`
	LdapAccess           string `json:"ldap_access"`
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
// @Param groupId body string false "group ID"
// @Param pageToken body string false "page Token"
// @Success 200  {object} []models.GitlabProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	groupId, ok := input.Query["groupId"]
	if !ok || len(groupId) == 0 {
		groupId = []string{""}
	}

	pageToken, ok := input.Query["pageToken"]
	if !ok || len(pageToken) == 0 {
		pageToken = []string{""}
	}

	// get gid and pageData
	gid := groupId[0]
	pageData, err := GetPageDataFromPageToken(pageToken[0])
	if err != nil {
		return nil, errors.BadInput.New("failed to get paget token")
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	outputBody := &RemoteScopesOutput{}

	// list groups part
	if pageData.Tag == TypeGroup {
		query, err := GetQueryFromPageData(pageData)
		if err != nil {
			return nil, err
		}
		query.Set("top_level_only", "true")

		if gid == "" {
			res, err = apiClient.Get("groups", query, nil)
		} else {
			res, err = apiClient.Get(fmt.Sprintf("groups/%s/subgroups", gid), query, nil)
		}
		if err != nil {
			return nil, err
		}

		resBody := []GroupResponse{}
		err = api.UnmarshalResponse(res, &resBody)

		// append group to output
		for _, group := range resBody {
			child := RemoteScopesChild{
				Type: TypeGroup,
				Id:   strconv.Itoa(group.Id),
				// don't need to save group into data
				Data: nil,
			}

			// ignore not top_level
			if group.ParentId == nil {
				if gid != "" {
					continue
				}
			} else {
				if strconv.Itoa(*group.ParentId) != gid {
					continue
				}
			}

			// ignore self
			if gid == child.Id {
				continue
			}

			child.ParentId = &gid
			if *child.ParentId == "" {
				child.ParentId = nil
			}

			outputBody.Children = append(outputBody.Children, child)
		}

		// check groups count
		if err != nil {
			return nil, err
		}
		if len(resBody) < pageData.PerPage {
			pageData.Tag = TypeProject
			pageData.Page = 1
			pageData.PerPage = pageData.PerPage - len(resBody)
		}
	}

	// list projects part
	if pageData.Tag == TypeProject {
		query, err := GetQueryFromPageData(pageData)
		if err != nil {
			return nil, err
		}
		if gid == "" {
			res, err = apiClient.Get(fmt.Sprintf("users/%d/projects", apiClient.GetData(models.GitlabApiClientData_UserId)), query, nil)
		} else {
			query.Set("with_shared", "false")
			res, err = apiClient.Get(fmt.Sprintf("/groups/%s/projects", gid), query, nil)
		}
		if err != nil {
			return nil, err
		}

		resBody := []tasks.GitlabApiProject{}
		err = api.UnmarshalResponse(res, &resBody)

		// append project to output
		for _, project := range resBody {
			child := RemoteScopesChild{
				Type: TypeProject,
				Id:   strconv.Itoa(project.CreatorId),
				Data: tasks.ConvertProject(&project),
			}
			child.ParentId = &gid
			if *child.ParentId == "" {
				child.ParentId = nil
			}

			outputBody.Children = append(outputBody.Children, child)
		}

		// check project count
		if err != nil {
			return nil, err
		}
		if len(resBody) < pageData.PerPage {
			pageData = nil
		}
	}

	// get the next page token
	outputBody.NextPageToken = ""
	if pageData != nil {
		pageData.Page += 1
		pageData.PerPage = GitlabRemoteScopesPerPage

		outputBody.NextPageToken, err = GetPageTokenFromPageData(pageData)
		if err != nil {
			return nil, err
		}
	}

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search body string false "group ID"
// @Param page body int false "page number"
// @Param pageSize body int false "page size per page"
// @Success 200  {object} []models.GitlabProject
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

	page, ok := input.Query["page"]
	if !ok || len(page) == 0 {
		page = []string{""}
	}

	p, err1 := strconv.Atoi(page[0])
	if err1 != nil {
		return nil, errors.BadInput.Wrap(err1, fmt.Sprintf("failed to Atoi page:%s", page[0]))
	}

	pageSize, ok := input.Query["pageSize"]
	if !ok || len(pageSize) == 0 {
		pageSize = []string{""}
	}
	ps, err1 := strconv.Atoi(pageSize[0])
	if err1 != nil {
		return nil, errors.BadInput.Wrap(err1, fmt.Sprintf("failed to Atoi pageSize:%s", pageSize[0]))
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

	// set projects return
	projects := []models.GitlabProject{}
	for _, project := range resBody {
		projects = append(projects, *tasks.ConvertProject(&project))
	}

	return &plugin.ApiResourceOutput{Body: projects, Status: http.StatusOK}, nil
}

func GetPageTokenFromPageData(pageData *PageData) (string, errors.Error) {
	// Marshal json
	pageTokenDecode, err := json.Marshal(pageData)
	if err != nil {
		return "", errors.Default.Wrap(err, fmt.Sprintf("Marshal pageToken failed %+v", pageData))
	}

	// Encode pageToken Base64
	return base64.StdEncoding.EncodeToString(pageTokenDecode), nil
}

func GetPageDataFromPageToken(pageToken string) (*PageData, errors.Error) {
	if pageToken == "" {
		return &PageData{
			Page:    1,
			PerPage: GitlabRemoteScopesPerPage,
			Tag:     "group",
		}, nil
	}

	// Decode pageToken Base64
	pageTokenDecode, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("decode pageToken failed %s", pageToken))
	}
	// Unmarshal json
	pt := &PageData{}
	err = json.Unmarshal(pageTokenDecode, pt)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("json Unmarshal pageTokenDecode failed %s", pageTokenDecode))
	}

	return pt, nil
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
