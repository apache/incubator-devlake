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

package models

import (
	"context"
	"fmt"
	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

var _ plugin.ApiConnectionForRemote[GroupResponse, BitbucketApiRepo] = (*BitbucketConnection)(nil)
var _ plugin.ApiGroup = (*GroupResponse)(nil)

// BitbucketConn holds the essential information to connect to the Bitbucket API
type BitbucketConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.BasicAuth      `mapstructure:",squash"`
}

// BitbucketConnection holds BitbucketConn plus ID/Name for database storage
type BitbucketConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	BitbucketConn      `mapstructure:",squash"`
}

func (BitbucketConnection) TableName() string {
	return "_tool_bitbucket_connections"
}

func (g BitbucketConnection) GetGroup(basicRes context2.BasicRes, gid string, query url.Values) ([]GroupResponse, errors.Error) {
	if gid != "" {
		return nil, nil
	}
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &g)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
	}
	var res *http.Response
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

	return resBody.Values, err
}

func (g BitbucketConnection) GetScope(basicRes context2.BasicRes, gid string, query url.Values) ([]BitbucketApiRepo, errors.Error) {
	if gid == "" {
		return nil, nil
	}
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &g)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
	}
	var res *http.Response
	query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.username,values.created_on,values.updated_on,values.links.clone,values.links.self,pagelen,page,size")
	// list projects part
	res, err = apiClient.Get(fmt.Sprintf("/repositories/%s", gid), query, nil)
	if err != nil {
		return nil, err
	}
	var resBody ReposResponse
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	return resBody.Values, err
}
