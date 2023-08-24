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
	"github.com/apache/incubator-devlake/core/log"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
)

type RemoteScopesChild struct {
	Type     string      `json:"type"`
	ParentId *string     `json:"parentId"`
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	FullName string      `json:"fullName"`
	Data     interface{} `json:"data"`
}

type RemoteQueryData struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	CustomInfo string `json:"custom"`
	Tag        string `json:"tag"`
	Search     []string
}

type FirstPageTokenOutput struct {
	PageToken string `json:"pageToken"`
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

// RemoteApiHelper is used to write the CURD of connection
type RemoteApiHelper[Conn plugin.ApiConnection, Scope plugin.ToolLayerScope, ApiScope plugin.ApiScope, Group plugin.ApiGroup] struct {
	basicRes   context.BasicRes
	validator  *validator.Validate
	connHelper *ConnectionApiHelper
	logger     log.Logger
}

// NewRemoteHelper creates a ScopeHelper for connection management
func NewRemoteHelper[Conn plugin.ApiConnection, Scope plugin.ToolLayerScope, ApiScope plugin.ApiScope, Group plugin.ApiGroup](
	basicRes context.BasicRes,
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
		logger:     basicRes.GetLogger(),
	}
}

type NoRemoteGroupResponse struct {
}

func (NoRemoteGroupResponse) GroupId() string {
	return ""
}

func (NoRemoteGroupResponse) GroupName() string {
	return ""
}

type BaseRemoteGroupResponse struct {
	Id   string
	Name string
}

func (g BaseRemoteGroupResponse) GroupId() string {
	return g.Id
}

func (g BaseRemoteGroupResponse) GroupName() string {
	return g.Name
}

const remoteScopesPerPage int = 100
const TypeScope string = "scope" // scope, sometimes we call it project. But scope is a more standard noun.
const TypeGroup string = "group" // group is just like a directory or a folder, that holds some scopes.

// PrepareFirstPageToken prepares the first page token
func (r *RemoteApiHelper[Conn, Scope, ApiScope, Group]) PrepareFirstPageToken(customInfo string) (*plugin.ApiResourceOutput, errors.Error) {
	outputBody := &FirstPageTokenOutput{}
	pageToken, err := getPageTokenFromPageData(&RemoteQueryData{
		Page:       1,
		PerPage:    remoteScopesPerPage,
		CustomInfo: customInfo,
		Tag:        TypeGroup,
	})
	if err != nil {
		return nil, err
	}
	outputBody.PageToken = pageToken
	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

// GetScopesFromRemote gets the scopes from api
func (r *RemoteApiHelper[Conn, Scope, ApiScope, Group]) GetScopesFromRemote(
	input *plugin.ApiResourceInput,
	getGroup func(basicRes context.BasicRes, gid string, queryData *RemoteQueryData, connection Conn) ([]Group, errors.Error),
	getScope func(basicRes context.BasicRes, gid string, queryData *RemoteQueryData, connection Conn) ([]ApiScope, errors.Error),
) (*plugin.ApiResourceOutput, errors.Error) {

	connectionId, err := errors.Convert01(strconv.ParseUint(input.Params["connectionId"], 10, 64))
	if err != nil || connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var connection Conn
	err = r.connHelper.First(&connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get groupId and pageData
	groupId := input.Query.Get("groupId")
	pageToken := input.Query.Get("pageToken")

	queryData, err := getPageDataFromPageToken(pageToken)
	if err != nil {
		return nil, errors.BadInput.New("failed to get page token")
	}
	outputBody := &RemoteScopesOutput{}

	// list groups part
	if queryData.Tag == TypeGroup {
		var resBody []Group
		if getGroup != nil {
			resBody, err = getGroup(r.basicRes, groupId, queryData, connection)
		}
		if err != nil {
			return nil, err
		}
		// if len(resBody) == 0, will skip the following steps, this will happen in some plugins which don't have group
		// append group to output
		for _, group := range resBody {
			child := RemoteScopesChild{
				Type: TypeGroup,
				Id:   group.GroupId(),
				Name: group.GroupName(),
				// don't need to save group into data
				Data: nil,
			}
			child.ParentId = &groupId
			if *child.ParentId == "" {
				child.ParentId = nil
			}
			outputBody.Children = append(outputBody.Children, child)
		}
		// check groups count
		if len(resBody) < queryData.PerPage {
			queryData.Tag = TypeScope
			queryData.Page = 1
			queryData.PerPage = queryData.PerPage - len(resBody)
		}
	}

	// list projects part
	if queryData.Tag == TypeScope && getScope != nil {
		var resBody []ApiScope
		resBody, err = getScope(r.basicRes, groupId, queryData, connection)
		if err != nil {
			return nil, err
		}

		// append project to output
		for _, project := range resBody {
			scope := project.ConvertApiScope()
			child := RemoteScopesChild{
				Type:     TypeScope,
				Id:       scope.ScopeId(),
				Name:     scope.ScopeName(),
				FullName: scope.ScopeFullName(),
				Data:     &scope,
			}
			child.ParentId = &groupId
			if *child.ParentId == "" {
				child.ParentId = nil
			}
			outputBody.Children = append(outputBody.Children, child)
		}

		// check scopes count
		if len(resBody) < queryData.PerPage {
			queryData = nil
		}
	}

	// get the next page token
	outputBody.NextPageToken = ""
	if queryData != nil {
		queryData.Page += 1
		outputBody.NextPageToken, err = getPageTokenFromPageData(queryData)
		if err != nil {
			return nil, err
		}
	}
	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

func (r *RemoteApiHelper[Conn, Scope, ApiScope, Group]) SearchRemoteScopes(input *plugin.ApiResourceInput,
	searchScope func(basicRes context.BasicRes, queryData *RemoteQueryData, connection Conn) ([]ApiScope, errors.Error),
) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, err := errors.Convert01(strconv.ParseUint(input.Params["connectionId"], 10, 64))
	if err != nil || connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	var connection Conn
	err = r.connHelper.First(&connection, input.Params)
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
		ps = remoteScopesPerPage
	} else {
		ps, err1 = strconv.Atoi(pageSize[0])
		if err1 != nil {
			return nil, errors.BadInput.Wrap(err1, fmt.Sprintf("failed to Atoi pageSize:%s", pageSize[0]))
		}
	}

	queryData := &RemoteQueryData{
		Page:    p,
		PerPage: ps,
		Search:  search,
	}

	var resBody []ApiScope
	resBody, err = searchScope(r.basicRes, queryData, connection)
	if err != nil {
		return nil, err
	}

	outputBody := &SearchRemoteScopesOutput{}

	// append project to output
	for _, project := range resBody {
		scope := project.ConvertApiScope()
		child := RemoteScopesChild{
			Type:     TypeScope,
			Id:       scope.ScopeId(),
			ParentId: nil,
			Name:     scope.ScopeName(),
			FullName: scope.ScopeFullName(),
			Data:     scope,
		}

		outputBody.Children = append(outputBody.Children, child)
	}

	outputBody.Page = p
	outputBody.PageSize = ps

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

func getPageTokenFromPageData(pageData *RemoteQueryData) (string, errors.Error) {
	// Marshal json
	pageTokenDecode, err := json.Marshal(pageData)
	if err != nil {
		return "", errors.Default.Wrap(err, fmt.Sprintf("Marshal pageToken failed %+v", pageData))
	}

	// Encode pageToken Base64
	return base64.StdEncoding.EncodeToString(pageTokenDecode), nil
}

func getPageDataFromPageToken(pageToken string) (*RemoteQueryData, errors.Error) {
	if pageToken == "" {
		return &RemoteQueryData{
			Page:    1,
			PerPage: remoteScopesPerPage,
			Tag:     TypeGroup,
		}, nil
	}

	// Decode pageToken Base64
	pageTokenDecode, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("decode pageToken failed %s", pageToken))
	}
	// Unmarshal json
	pt := &RemoteQueryData{}
	err = json.Unmarshal(pageTokenDecode, pt)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("json Unmarshal pageTokenDecode failed %s", pageTokenDecode))
	}

	return pt, nil
}
