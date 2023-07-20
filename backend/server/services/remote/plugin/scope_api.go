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

	"github.com/apache/incubator-devlake/server/services/remote/models"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type request struct {
	Data []map[string]any `json:"data"`
}

func (pa *pluginAPI) PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var scopes request
	err := errors.Convert(mapstructure.Decode(input.Body, &scopes))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding scope error")
	}
	var slice []*models.DynamicScopeModel
	for _, scope := range scopes.Data {
		obj := models.NewDynamicScopeModel(pa.scopeType)
		err = models.MapTo(scope, obj)
		if err != nil {
			return nil, err
		}
		slice = append(slice, obj)
	}
	apiScopes, err := pa.scopeHelper.PutScopes(input, slice)
	if err != nil {
		return nil, err
	}
	response, err := convertScopeResponse(apiScopes...)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	apiScopes, err := pa.scopeHelper.UpdateScope(input)
	if err != nil {
		return nil, err
	}
	response, err := convertScopeResponse(apiScopes)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: response[0], Status: http.StatusOK}, nil
}

func (pa *pluginAPI) ListScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopes, err := pa.scopeHelper.GetScopes(input)
	if err != nil {
		return nil, err
	}
	response, err := convertScopeResponse(scopes...)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scope, err := pa.scopeHelper.GetScope(input)
	if err != nil {
		return nil, err
	}
	response, err := convertScopeResponse(scope)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: response[0], Status: http.StatusOK}, nil
}

func (pa *pluginAPI) DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	refs, err := pa.scopeHelper.DeleteScope(input)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: refs, Status: err.GetType().GetHttpCode()}, nil
	}
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

// convertScopeResponse adapt the "remote" scopes to a serializable api.ScopeRes. This code is needed because squashed mapstructure don't work
// with dynamic/runtime structs used by remote plugins
func convertScopeResponse(scopes ...*api.ScopeRes[models.DynamicScopeModel, models.RemoteScopeConfig]) ([]map[string]any, errors.Error) {
	responses := make([]map[string]any, len(scopes))
	for i, scope := range scopes {
		resMap := map[string]any{}
		err := models.MapTo(scope.ScopeResDoc, &resMap)
		if err != nil {
			return nil, err
		}
		scopeMap := map[string]map[string]any{}
		err = models.MapTo(scope.Scope, &scopeMap)
		if err != nil {
			return nil, err
		}
		for k, v := range scopeMap["DynamicTabler"] {
			resMap[k] = v
		}
		responses[i] = resMap
	}
	return responses, nil
}
