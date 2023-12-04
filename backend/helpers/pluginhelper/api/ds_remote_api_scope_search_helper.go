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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/helpers/utils"
)

type DsSearchRemoteScopes[C plugin.ToolLayerApiConnection, S plugin.ToolLayerScope] func(
	apiClient plugin.ApiClient, params *models.DsRemoteApiScopeSearchParams) (children []models.DsRemoteApiScopeListEntry[S], err errors.Error)

type DsRemoteApiScopeSearchHelper[C plugin.ToolLayerApiConnection, S plugin.ToolLayerScope] struct {
	*DsRemoteApiProxyHelper[C]
	searchRemoteScopes DsSearchRemoteScopes[C, S]
}

func NewDsRemoteApiScopeSearchHelper[
	C plugin.ToolLayerApiConnection,
	S plugin.ToolLayerScope,
](
	rap *DsRemoteApiProxyHelper[C],
	searchRemoteScopes DsSearchRemoteScopes[C, S],
) *DsRemoteApiScopeSearchHelper[C, S] {
	return &DsRemoteApiScopeSearchHelper[C, S]{
		DsRemoteApiProxyHelper: rap,
		searchRemoteScopes:     searchRemoteScopes,
	}
}

// SearchRemoteScopes searches remote scopes
func (rss *DsRemoteApiScopeSearchHelper[C, S]) Get(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	_, apiClient, err := rss.prepare(input)
	if err != nil {
		return nil, err
	}
	params := &models.DsRemoteApiScopeSearchParams{
		Page:     1,
		PageSize: 50,
	}
	err = utils.DecodeMapStruct(input.Query, params, true)
	if err != nil {
		return nil, err
	}
	if e := vld.Struct(params); e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid params")
	}
	children, err := rss.searchRemoteScopes(apiClient, params)
	if err != nil {
		return nil, err
	}
	if children == nil {
		children = []models.DsRemoteApiScopeListEntry[S]{}
	}
	// the config-ui is expecting the parent id to be null
	for i := range children {
		children[i].ParentId = nil
	}
	return &plugin.ApiResourceOutput{
		Body: map[string]interface{}{
			"children": children,
			"page":     params.Page,
			"pageSize": params.PageSize,
		},
	}, nil
}
