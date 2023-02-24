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

package plugin

import (
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ScopeItem struct {
	ScopeId              string `json:"scopeId"`
	ScopeName            string `json:"scopeName"`
	ConnectionId         uint64 `json:"connectionId"`
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty"`
}

// DTO that includes the transformation rule name
type apiScopeResponse struct {
	Scope                  ScopeItem
	TransformationRuleName string `json:"transformationRuleId,omitempty"`
}

// Why a batch PUT?
type request struct {
	Data []*ScopeItem `json:"data"`
}

func (pa *pluginAPI) PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	var scopes request
	err := errors.Convert(mapstructure.Decode(input.Body, &scopes))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding scope error")
	}

	keeper := make(map[string]struct{})
	for _, scope := range scopes.Data {
		if _, ok := keeper[scope.ScopeId]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[scope.ScopeId] = struct{}{}
		}
		scope.ConnectionId = connectionId

		err = verifyScope(scope)
		if err != nil {
			return nil, err
		}
	}

	err = basicRes.GetDal().CreateOrUpdate(scopes.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving scope")
	}

	return &plugin.ApiResourceOutput{Body: scopes.Data, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, scopeId := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	db := basicRes.GetDal()
	scope := ScopeItem{}
	err := db.First(&scope, dal.Where("connection_id = ? AND scope_id = ?", connectionId, scopeId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "scope not found")
	}

	err = api.DecodeMapStruct(input.Body, &scope)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch scope error")
	}

	err = verifyScope(&scope)
	if err != nil {
		return nil, err
	}

	err = db.Update(&scope)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving scope")
	}
	return &plugin.ApiResourceOutput{Body: scope, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) ListScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var scopes []ScopeItem
	connectionId, _ := extractParam(input.Params)

	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")

	if limit > 100 {
		return nil, errors.BadInput.New("Page limit cannot exceed 100")
	}

	db := basicRes.GetDal()
	err := db.All(&scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}

	var ruleIds []uint64
	for _, scope := range scopes {
		if scope.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, scope.TransformationRuleId)
		}
	}

	var txRuleId2Name []struct {
		id   uint64
		name string
	}
	if len(ruleIds) > 0 {
		err = db.All(&txRuleId2Name,
			dal.Select("id, name"),
			dal.From(pa.txRuleType.TableName()),
			dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}

	names := make(map[uint64]string)
	for _, r := range txRuleId2Name {
		names[r.id] = r.name
	}

	var apiScopes []apiScopeResponse
	for _, scope := range scopes {
		txRuleName := names[scope.TransformationRuleId]
		scopeRes := apiScopeResponse{
			Scope:                  scope,
			TransformationRuleName: txRuleName,
		}
		apiScopes = append(apiScopes, scopeRes)
	}

	return &plugin.ApiResourceOutput{Body: apiScopes, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var scope ScopeItem
	connectionId, scopeId := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}

	db := basicRes.GetDal()
	err := db.First(&scope, dal.Where("connection_id = ? AND scope_id = ?", connectionId, scopeId))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}

	var ruleName string
	if scope.TransformationRuleId > 0 {
		err = db.First(&ruleName, dal.Select("name"), dal.From(pa.txRuleType.TableName()), dal.Where("id = ?", scope.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}

	return &plugin.ApiResourceOutput{Body: apiScopeResponse{scope, ruleName}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	scopeId := params["scopeId"]
	return connectionId, scopeId
}

func verifyScope(scope *ScopeItem) errors.Error {
	if scope.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}

	if scope.ScopeId == "" {
		return errors.BadInput.New("invalid scope ID")
	}

	return nil
}
