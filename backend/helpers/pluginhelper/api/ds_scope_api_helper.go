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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	serviceHelper "github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/server/api/shared"
)

// for documentation purposes
type ScopeRefDoc = serviceHelper.BlueprintProjectPairs
type PutScopesReqBody[T any] struct {
	Data []*T `json:"data"`
}
type ScopeDetail[S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	Scope       S                   `json:"scope"`
	ScopeConfig *SC                 `json:"scopeConfig,omitempty"`
	Blueprints  []*models.Blueprint `json:"blueprints,omitempty"`
}

type DsAnyScopeApiHelper struct {
	*AnyModelApiHelper
	*srvhelper.AnyScopeSrvHelper
}

func NewDsAnyScopeApiHelper(
	basicRes context.BasicRes,
	srvHelper *srvhelper.AnyScopeSrvHelper,
) *DsAnyScopeApiHelper {
	return &DsAnyScopeApiHelper{
		AnyModelApiHelper: NewAnyModelApiHelper(basicRes, srvHelper.AnyModelSrvHelper, []string{"connectionId", "scopeId"}, nil),
		AnyScopeSrvHelper: srvHelper,
	}
}

func (scopeApi *DsAnyScopeApiHelper) GetPage(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	pagination, err := parsePagination[srvhelper.ScopePagination](input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to decode pathvars into pagination")
	}
	scopes, count, err := scopeApi.GetScopesPageAny(pagination)
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

func (scopeApi *DsAnyScopeApiHelper) GetScopeDetail(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	pkv, err := scopeApi.ExtractPkValues(input)
	if err != nil {
		return nil, err
	}
	scopeDetail, err := scopeApi.GetScopeDetailAny(input.Query.Get("blueprints") == "true", pkv...)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: scopeDetail,
	}, nil
}

func (scopeApi *DsAnyScopeApiHelper) GetScopeLatestSyncState(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	pkv, err := scopeApi.ExtractPkValues(input)
	if err != nil {
		return nil, err
	}
	scopeLatestSyncStates, err := scopeApi.GetScopeLatestSyncStateAny(pkv...)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: scopeLatestSyncStates,
	}, nil
}

func (scopeApi *DsAnyScopeApiHelper) PutMultiple(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
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
	return scopeApi.PutMultipleCb(input, func(m any) errors.Error {
		ok := setRawDataOrigin(m, common.RawDataOrigin{
			RawDataTable:  fmt.Sprintf("_raw_%s_scopes", scopeApi.GetPluginName()),
			RawDataParams: plugin.MarshalScopeParams(scopeApi.GetScopeParams(m)),
		})
		if !ok {
			panic("set raw data origin failed")
		}
		return nil
	})
}

func (scopeApi *DsAnyScopeApiHelper) Delete(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scope, err := scopeApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	// time.Sleep(1 * time.Minute) # uncomment this line if you were to verify pipelines get blocked while deleting data
	// check referencing blueprints
	refs, err := scopeApi.DeleteScopeAny(scope, input.Query.Get("delete_data_only") == "true")
	if err != nil {
		return &plugin.ApiResourceOutput{Body: &shared.ApiBody{
			Success: false,
			Message: err.Error(),
			Data:    refs,
		}, Status: err.GetType().GetHttpCode()}, err
	}
	return &plugin.ApiResourceOutput{
		Body: scope,
	}, nil
}

type DsScopeApiHelper[S dal.Tabler] struct {
	*DsAnyScopeApiHelper
	*ModelApiHelper[S]
}

func NewDsScopeApiHelper[S dal.Tabler](
	anyScopeApiHelper *DsAnyScopeApiHelper,
) *DsScopeApiHelper[S] {
	return &DsScopeApiHelper[S]{
		DsAnyScopeApiHelper: anyScopeApiHelper,
		ModelApiHelper:      NewModelApiHelper[S](anyScopeApiHelper.AnyModelApiHelper),
	}
}
