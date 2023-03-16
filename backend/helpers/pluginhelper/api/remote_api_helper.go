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
	"encoding/base64"
	"encoding/json"
	"fmt"
	coreContext "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
)

// RemoteApiHelper is used to write the CURD of connection
type RemoteApiHelper[Conn plugin.ApiConnectionForRemote[Group, ApiScope], Scope plugin.ToolLayerScope, ApiScope plugin.ApiScope, Group plugin.ApiGroup] struct {
	basicRes   coreContext.BasicRes
	validator  *validator.Validate
	connHelper *ConnectionApiHelper
}

// NewRemoteHelper creates a ScopeHelper for connection management
func NewRemoteHelper[Conn plugin.ApiConnectionForRemote[Group, ApiScope], Scope plugin.ToolLayerScope, ApiScope plugin.ApiScope, Group plugin.ApiGroup](
	basicRes coreContext.BasicRes,
	vld *validator.Validate,
	connHelper *ConnectionApiHelper,
) *RemoteApiHelper[Conn, Scope, ApiScope, Group] {
	if vld == nil {
		vld = validator.New()
	}
	if connHelper == nil {
		return nil
	}
	return &RemoteApiHelper[Conn, Scope, ApiScope, Group]{
		basicRes:   basicRes,
		validator:  vld,
		connHelper: connHelper,
	}
}

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

const remoteScopesPerPage int = 100
const TypeProject string = "scope"
const TypeGroup string = "group"

func (r *RemoteApiHelper[Conn, Scope, ApiScope, Group]) GetScopesFromRemote(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := ExtractFromReqParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	var connection Conn
	err := r.connHelper.First(&connection, input.Params)
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
	pageData, err := getPageDataFromPageToken(pageToken[0])
	if err != nil {
		return nil, errors.BadInput.New("failed to get paget token")
	}

	outputBody := &RemoteScopesOutput{}

	// list groups part
	if pageData.Tag == TypeGroup {
		query, err := getQueryFromPageData(pageData)
		if err != nil {
			return nil, err
		}
		var resBody []Group
		resBody, err = connection.GetGroup(r.basicRes, gid, query)
		if err != nil {
			return nil, err
		}

		// append group to output
		for _, group := range resBody {
			child := RemoteScopesChild{
				Type: TypeGroup,
				Id:   group.GroupId(),
				Name: group.GroupName(),
				// don't need to save group into data
				Data: nil,
			}
			child.ParentId = &gid
			if *child.ParentId == "" {
				child.ParentId = nil
			}
			outputBody.Children = append(outputBody.Children, child)
		}
		// check groups count
		if len(resBody) < pageData.PerPage {
			pageData.Tag = TypeProject
			pageData.Page = 1
			pageData.PerPage = pageData.PerPage - len(resBody)
		}
	}

	// list projects part
	if pageData.Tag == TypeProject {
		query, err := getQueryFromPageData(pageData)
		if err != nil {
			return nil, err
		}
		var resBody []ApiScope

		resBody, err = connection.GetScope(r.basicRes, gid, query)
		if err != nil {
			return nil, err
		}

		// append project to output
		for _, project := range resBody {
			scope := project.ConvertApiScope()
			child := RemoteScopesChild{
				Type: TypeProject,
				Id:   scope.ScopeId(),
				Name: scope.ScopeName(),
				Data: &scope,
			}
			child.ParentId = &gid
			if *child.ParentId == "" {
				child.ParentId = nil
			}

			outputBody.Children = append(outputBody.Children, child)
		}

		// check project count
		if len(resBody) < pageData.PerPage {
			pageData = nil
		}
	}

	// get the next page token
	outputBody.NextPageToken = ""
	if pageData != nil {
		pageData.Page += 1
		pageData.PerPage = remoteScopesPerPage

		outputBody.NextPageToken, err = getPageTokenFromPageData(pageData)
		if err != nil {
			return nil, err
		}
	}

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

func getPageTokenFromPageData(pageData *PageData) (string, errors.Error) {
	// Marshal json
	pageTokenDecode, err := json.Marshal(pageData)
	if err != nil {
		return "", errors.Default.Wrap(err, fmt.Sprintf("Marshal pageToken failed %+v", pageData))
	}

	// Encode pageToken Base64
	return base64.StdEncoding.EncodeToString(pageTokenDecode), nil
}

func getPageDataFromPageToken(pageToken string) (*PageData, errors.Error) {
	if pageToken == "" {
		return &PageData{
			Page:    1,
			PerPage: remoteScopesPerPage,
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

func getQueryFromPageData(pageData *PageData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", pageData.Page))
	query.Set("per_page", fmt.Sprintf("%v", pageData.PerPage))
	return query, nil
}

func GetQueryForSearchProject(search string, page int, perPage int) (url.Values, errors.Error) {
	query, err := getQueryFromPageData(&PageData{Page: page, PerPage: perPage})
	if err != nil {
		return nil, err
	}
	query.Set("search", search)
	query.Set("scope", "projects")

	return query, nil
}
