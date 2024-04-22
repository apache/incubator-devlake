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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type request struct {
	Data []map[string]any `json:"data"`
}

func (pa *pluginAPI) PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.PutMultiple(input)
}

func (pa *pluginAPI) UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.Patch(input)
}

func (pa *pluginAPI) ListScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.GetPage(input)
}

func (pa *pluginAPI) GetScopeDispatcher(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeIdWithSuffix := strings.TrimLeft(input.Params["scopeId"], "/")
	if strings.HasSuffix(scopeIdWithSuffix, "/latest-sync-state") {
		input.Params["scopeId"] = strings.TrimSuffix(scopeIdWithSuffix, "/latest-sync-state")
		return pa.GetScopeLatestSyncState(input)
	}
	return pa.GetScope(input)
}

func (pa *pluginAPI) GetScopeLatestSyncState(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.GetScopeLatestSyncState(input)
}

func (pa *pluginAPI) GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.GetDetail(input)
}

func (pa *pluginAPI) DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return pa.dsHelper.ScopeApi.Delete(input)
}
