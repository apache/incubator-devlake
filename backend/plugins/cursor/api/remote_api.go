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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helperapi "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/helpers/utils"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
)

// CursorRemotePagination is a placeholder for remote scope pagination.
type CursorRemotePagination struct {
	Page int `json:"page"`
}

func listCursorRemoteScopes(
	_ *models.CursorConnection,
	_ plugin.ApiClient,
	_ string,
	_ CursorRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.CursorScope],
	nextPage *CursorRemotePagination,
	err errors.Error,
) {
	children = append(children, makeCursorRemoteScopeEntry())
	return children, nil, nil
}

func makeCursorRemoteScopeEntry() dsmodels.DsRemoteApiScopeListEntry[models.CursorScope] {
	return dsmodels.DsRemoteApiScopeListEntry[models.CursorScope]{
		Type:     helperapi.RAS_ENTRY_TYPE_SCOPE,
		Id:       models.DefaultScopeID,
		Name:     "Team",
		FullName: "Cursor Team",
		Data: &models.CursorScope{
			Id:       models.DefaultScopeID,
			Name:     "Team",
			FullName: "Cursor Team",
		},
	}
}

// RemoteScopes lists all available scopes for this connection.
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes searches scopes for this connection.
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	params := &dsmodels.DsRemoteApiScopeSearchParams{
		Page:     1,
		PageSize: 50,
	}
	if err := utils.DecodeMapStruct(input.Query, params, true); err != nil {
		return nil, err
	}
	if err := errors.Convert(vld.Struct(params)); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid params")
	}

	children := []dsmodels.DsRemoteApiScopeListEntry[models.CursorScope]{}
	searchLower := strings.ToLower(strings.TrimSpace(params.Search))
	if searchLower == "" || searchLower == "team" || searchLower == "cursor" {
		children = append(children, makeCursorRemoteScopeEntry())
	}

	return &plugin.ApiResourceOutput{
		Body: map[string]interface{}{
			"children": children,
			"page":     params.Page,
			"pageSize": params.PageSize,
		},
	}, nil
}
