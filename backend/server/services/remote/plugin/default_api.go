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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

type pluginAPI struct {
	invoker  bridge.Invoker
	dsHelper *api.DsAnyHelper
}

func GetDefaultAPI(
	invoker bridge.Invoker,
	dsHelper *api.DsAnyHelper,
) map[string]map[string]plugin.ApiResourceHandler {
	papi := &pluginAPI{
		invoker:  invoker,
		dsHelper: dsHelper,
	}
	resources := map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": papi.TestConnection,
		},
		"connections": {
			"POST": papi.PostConnections,
			"GET":  papi.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    papi.GetConnection,
			"PATCH":  papi.PatchConnection,
			"DELETE": papi.DeleteConnection,
		},
		"connections/:connectionId/test": {
			"POST": papi.TestExistingConnection,
		},
		"connections/:connectionId/scopes": {
			"PUT": papi.PutScope,
			"GET": papi.ListScopes,
		},
		// Use `*` to match scopeId with `/` in it
		"connections/:connectionId/scopes/*scopeId": {
			// Behind 'GetScopeDispatcher', there are two paths so far:
			// GetScopeLatestSyncState "connections/:connectionId/scopes/:scopeId/latest-sync-state"
			// GetScope "connections/:connectionId/scopes/:scopeId"
			// Because there may be slash in scopeId, so we handle it manually.
			"GET":    papi.GetScopeDispatcher,
			"PATCH":  papi.UpdateScope,
			"DELETE": papi.DeleteScope,
		},
		"connections/:connectionId/scope-configs": {
			"POST": papi.PostScopeConfigs,
			"GET":  papi.ListScopeConfigs,
		},
		"connections/:connectionId/scope-configs/:id": {
			"GET":    papi.GetScopeConfig,
			"PATCH":  papi.PatchScopeConfig,
			"DELETE": papi.DeleteScopeConfig,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": papi.GetRemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": papi.SearchRemoteScopes,
		},
		"scope-config/:id/projects": {
			"GET": papi.GetProjectsByScopeConfig,
		},
	}
	return resources
}
