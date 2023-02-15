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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type WorkspaceResponse struct {
	Pagelen int `json:"pagelen"`
	Page    int `json:"page"`
	Size    int `json:"size"`
	Values  []struct {
		//Type       string `json:"type"`
		//Permission string `json:"permission"`
		//LastAccessed time.Time `json:"last_accessed"`
		//AddedOn      time.Time `json:"added_on"`
		Workspace WorkspaceItem `json:"workspace"`
	} `json:"values"`
}

type WorkspaceItem struct {
	//Type string `json:"type"`
	//Uuid string `json:"uuid"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type ReposResponse struct {
	Pagelen int                      `json:"pagelen"`
	Page    int                      `json:"page"`
	Size    int                      `json:"size"`
	Values  []tasks.BitbucketApiRepo `json:"values"`
}

const RemoteScopesPerPage int = 100
const TypeScope string = "scope"
const TypeGroup string = "group"

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.BitbucketConnection{}
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
		return nil, errors.BadInput.New("failed to get page token")
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	query, err := GetQueryFromPageData(pageData)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	outputBody := &RemoteScopesOutput{}

	// list groups part
	if gid == "" {
		query.Set("sort", "workspace.slug")
		query.Set("fields", "values.workspace.slug,values.workspace.name,pagelen,page,size")
		res, err = apiClient.Get("/user/permissions/workspaces", query, nil)
		if err != nil {
			return nil, err
		}

		resBody := &WorkspaceResponse{}
		err = api.UnmarshalResponse(res, resBody)
		if err != nil {
			return nil, err
		}

		// append group to output
		for _, group := range resBody.Values {
			child := RemoteScopesChild{
				Type: TypeGroup,
				Id:   group.Workspace.Slug,
				Name: group.Workspace.Name,
				// don't need to save group into data
				Data: nil,
			}
			outputBody.Children = append(outputBody.Children, child)
		}

		// check groups count
		if resBody.Size < pageData.PerPage {
			pageData = nil
		}
	} else {
		query.Set("sort", "name")
		query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.username,values.created_on,values.updated_on,values.links.clone,values.links.self,pagelen,page,size")
		// list projects part
		res, err = apiClient.Get(fmt.Sprintf("/repositories/%s", gid), query, nil)
		if err != nil {
			return nil, err
		}

		resBody := &ReposResponse{}
		err = api.UnmarshalResponse(res, resBody)
		if err != nil {
			return nil, err
		}

		// append repo to output
		for _, repo := range resBody.Values {
			child := RemoteScopesChild{
				Type: TypeScope,
				Id:   repo.FullName,
				Name: repo.Name,
				Data: tasks.ConvertApiRepoToScope(&repo, connection.ID),
			}
			child.ParentId = &gid
			outputBody.Children = append(outputBody.Children, child)
		}

		// check repo count
		if resBody.Size < pageData.PerPage {
			pageData = nil
		}
	}

	// get the next page token
	outputBody.NextPageToken = ""
	if pageData != nil {
		pageData.Page += 1
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
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.BitbucketConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	search, ok := input.Query["search"]
	if !ok || len(search) == 0 {
		search = []string{""}
	}
	s := search[0]

	p := 1
	page, ok := input.Query["page"]
	if ok && len(page) != 0 {
		p, err = errors.Convert01(strconv.Atoi(page[0]))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("failed to Atoi page:%s", page[0]))
		}
	}

	ps := RemoteScopesPerPage
	pageSize, ok := input.Query["pageSize"]
	if ok && len(pageSize) != 0 {
		ps, err = errors.Convert01(strconv.Atoi(pageSize[0]))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("failed to Atoi pageSize:%s", pageSize[0]))
		}
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	// set query
	query, err := GetQueryFromPageData(&PageData{p, ps})
	if err != nil {
		return nil, err
	}

	// request search
	query.Set("sort", "name")
	query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.username,values.created_on,values.updated_on,values.links.clone,values.links.self,pagelen,page,size")

	gid := ``
	if strings.Contains(s, `/`) {
		gid = strings.Split(s, `/`)[0]
		s = strings.Split(s, `/`)[0]
	}
	query.Set("q", fmt.Sprintf(`name~"%s"`, s))
	// list repos part
	res, err := apiClient.Get(fmt.Sprintf("/repositories/%s", gid), query, nil)
	if err != nil {
		return nil, err
	}
	resBody := &ReposResponse{}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}

	// set repos return
	outputBody := &SearchRemoteScopesOutput{Children: []RemoteScopesChild{}}
	for _, repo := range resBody.Values {
		child := RemoteScopesChild{
			Type:     TypeScope,
			Id:       repo.FullName,
			ParentId: nil,
			Name:     repo.Name,
			Data:     tasks.ConvertApiRepoToScope(&repo, connection.ID),
		}
		outputBody.Children = append(outputBody.Children, child)
	}

	outputBody.Page = p
	outputBody.PageSize = ps

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
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
			PerPage: RemoteScopesPerPage,
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
	query.Set("pagelen", fmt.Sprintf("%v", pageData.PerPage))
	return query, nil
}
