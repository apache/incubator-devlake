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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type DsScopeApiHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*ModelApiHelper[S]
	*srvhelper.ScopeSrvHelper[C, S, SC]
}

func NewDsScopeApiHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig](
	basicRes context.BasicRes,
	srvHelper *srvhelper.ScopeSrvHelper[C, S, SC],
) *DsScopeApiHelper[C, S, SC] {
	return &DsScopeApiHelper[C, S, SC]{
		ModelApiHelper: NewModelApiHelper[S](basicRes, srvHelper.ModelSrvHelper, []string{"connectionId", "scopeId"}),
		ScopeSrvHelper: srvHelper,
	}
}

func (connApi *DsScopeApiHelper[C, S, SC]) GetPage(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	pagination, err := ParsePagination[srvhelper.ScopePagination](input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to decode pathvars into pagination")
	}
	scopes, count, err := connApi.ScopeSrvHelper.GetScopesPage(pagination)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: map[string]interface{}{
			"count":  count,
			"scopes": scopes,
		},
	}, nil
}

func (connApi *DsScopeApiHelper[C, S, SC]) PutMultiple(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// fix data[].connectionId
	connectionId, err := extractConnectionId(input)
	if err != nil {
		return nil, err
	}
	data, ok := input.Body["data"].([]interface{})
	if !ok {
		return nil, errors.BadInput.New("invalid data")
	}
	for _, row := range data {
		dict, ok := row.(map[string]interface{})
		if !ok {
			return nil, errors.BadInput.New("invalid data row")
		}
		dict["connectionId"] = connectionId
	}
	return connApi.ModelApiHelper.PutMultiple(input)
}

func (connApi *DsScopeApiHelper[C, S, SC]) Delete(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var scope *S
	scope, err := connApi.FindByPk(input)
	if err != nil {
		return nil, err
	}
	// time.Sleep(1 * time.Minute) # uncomment this line if you were to verify pipelines get blocked while deleting data
	// check referencing blueprints
	refs, err := connApi.ScopeSrvHelper.DeleteScope(scope, input.Query.Get("delete_data_only") == "true")
	if err != nil {
		return &plugin.ApiResourceOutput{Body: &shared.ApiBody{
			Success: false,
			Message: err.Error(),
			Data:    refs,
		}, Status: err.GetType().GetHttpCode()}, nil
	}
	return &plugin.ApiResourceOutput{
		Body: scope,
	}, nil
}
