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
	"github.com/apache/incubator-devlake/server/services/remote/models"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

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
	scopes, err := scopeHelper.GetScopeList(input)
	if err != nil {
		return nil, err
	}
	var response []api.ScopeRes[models.ScopeModel]
	for _, scope := range scopes {
		scopeModel := models.ScopeModel{}
		err = mapTo(scope.Scope, &scopeModel)
		if err != nil {
			return nil, err
		}
		response = append(response, api.ScopeRes[models.ScopeModel]{
			Scope:                  scopeModel,
			TransformationRuleName: scope.TransformationRuleName,
			Blueprints:             scope.Blueprints,
		})
	}
	return &plugin.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scope, err := scopeHelper.GetScope(input)
	if err != nil {
		return nil, err
	}
	scopeModel := models.ScopeModel{}
	err = mapTo(scope.Scope, &scopeModel)
	if err != nil {
		return nil, err
	}
	response := api.ScopeRes[models.ScopeModel]{
		Scope:                  scopeModel,
		TransformationRuleName: scope.TransformationRuleName,
		Blueprints:             scope.Blueprints,
	}
	return &plugin.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	err := scopeHelper.Delete(input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
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

func mapTo(x any, y any) errors.Error {
	b, err := json.Marshal(x)
	if err != nil {
		return errors.Convert(err)
	}
	if err = json.Unmarshal(b, y); err != nil {
		return errors.Convert(err)
	}
	return nil
}
