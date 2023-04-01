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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/server/services/remote/models"

	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// DTO that includes the transformation rule name
type apiScopeResponse struct {
	Scope                  any    `json:"-"`
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

// MarshalJSON make Scope display inline
func (r apiScopeResponse) MarshalJSON() ([]byte, error) {
	// encode scope to map
	scopeBytes, err := json.Marshal(r.Scope)
	if err != nil {
		return nil, err
	}
	var scopeMap map[string]interface{}
	err = json.Unmarshal(scopeBytes, &scopeMap)
	if err != nil {
		return nil, err
	}

	// encode other column (transformationRuleName) to map
	otherBytes, err := json.Marshal(struct {
		TransformationRuleName string `json:"transformationRuleName,omitempty"`
	}{
		TransformationRuleName: r.TransformationRuleName,
	})
	if err != nil {
		return nil, err
	}

	// merge the two maps
	var merged map[string]interface{}
	err = json.Unmarshal(otherBytes, &merged)
	if err != nil {
		return nil, err
	}
	for k, v := range scopeMap {
		merged[k] = v
	}

	// encode the merged map to JSON
	return json.Marshal(merged)
}

type request struct {
	Data []map[string]any `json:"data"`
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
	var createdScopes []any
	for _, scopeRaw := range scopes.Data {
		err = verifyScope(scopeRaw)
		if err != nil {
			return nil, err
		}
		scopeId := scopeRaw["id"].(string)
		if _, ok := keeper[scopeId]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[scopeId] = struct{}{}
		}
		scope := pa.scopeType.New()
		err = scope.From(&scopeRaw)
		if err != nil {
			return nil, err
		}
		// I don't know the reflection logic to do this in a batch...
		err = api.CallDB(basicRes.GetDal().CreateOrUpdate, scope)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error on saving scope")
		}
		createdScopes = append(createdScopes, scope.Unwrap())
	}

	return &plugin.ApiResourceOutput{Body: createdScopes, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, scopeId := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	db := basicRes.GetDal()
	scope := pa.scopeType.New()
	err := api.CallDB(db.First, scope, dal.Where("connection_id = ? AND id = ?", connectionId, scopeId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "scope not found")
	}
	err = verifyScope(input.Body)
	if err != nil {
		return nil, err
	}
	err = scope.From(&input.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch scope error")
	}
	err = api.CallDB(db.Update, scope)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving scope")
	}
	return &plugin.ApiResourceOutput{Body: scope.Unwrap(), Status: http.StatusOK}, nil
}

func (pa *pluginAPI) ListScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")

	if limit > 100 {
		return nil, errors.BadInput.New("Page limit cannot exceed 100")
	}
	db := basicRes.GetDal()
	scopes := pa.scopeType.NewSlice()
	err := api.CallDB(db.All, scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var scopeMap []map[string]any
	err = scopes.To(&scopeMap)
	if err != nil {
		return nil, err
	}
	if pa.txRuleType == nil {
		var apiScopes []apiScopeResponse
		for _, scope := range scopeMap {
			apiScopes = append(apiScopes, apiScopeResponse{Scope: scope})
		}
		return &plugin.ApiResourceOutput{Body: apiScopes, Status: http.StatusOK}, nil
	}
	var ruleIds []uint64
	for _, scopeModel := range scopeMap {
		if tid := uint64(scopeModel["transformationRuleId"].(float64)); tid > 0 {
			ruleIds = append(ruleIds, tid)
		}
	}
	rules := pa.txRuleType.NewSlice()
	if len(ruleIds) > 0 {
		err = api.CallDB(db.All, rules, dal.Select("id, name"),
			dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	var transformationModels []models.TransformationModel
	err = rules.To(&transformationModels)
	if err != nil {
		return nil, err
	}
	names := make(map[uint64]string)
	for _, t := range transformationModels {
		names[t.Id] = t.Name
	}
	var apiScopes []apiScopeResponse
	for _, scope := range scopeMap {
		txRuleName, ok := names[uint64(scope["transformationRuleId"].(float64))]
		if ok {
			scopeRes := apiScopeResponse{
				Scope:                  scope,
				TransformationRuleName: txRuleName,
			}
			apiScopes = append(apiScopes, scopeRes)
		}
	}

	return &plugin.ApiResourceOutput{Body: apiScopes, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, scopeId := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	if scopeId == `` {
		return nil, errors.BadInput.New("invalid scopeId")
	}
	rawScope := pa.scopeType.New()
	db := basicRes.GetDal()
	err := api.CallDB(db.First, rawScope, dal.Where("connection_id = ? AND id = ?", connectionId, scopeId))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	var scope models.ScopeModel
	err = rawScope.To(&scope)
	if err != nil {
		return nil, err
	}
	var rule models.TransformationModel
	if scope.TransformationRuleId > 0 {
		err = api.CallDB(db.First, &rule, dal.From(pa.txRuleType.TableName()), dal.Where("id = ?", scope.TransformationRuleId))
		if err != nil {
			return nil, errors.Default.Wrap(err, `no related transformationRule for scope`)
		}
	}
	return &plugin.ApiResourceOutput{Body: apiScopeResponse{rawScope.Unwrap(), rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	scopeId := params["scopeId"]
	return connectionId, scopeId
}

func verifyScope(scope map[string]any) errors.Error {
	if connectionId, ok := scope["connectionId"]; !ok || connectionId.(float64) == 0 {
		return errors.BadInput.New("invalid connectionId")
	}

	if scope["id"] == "" {
		return errors.BadInput.New("invalid scope ID")
	}

	return nil
}
