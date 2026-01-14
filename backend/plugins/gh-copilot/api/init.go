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
	"github.com/go-playground/validator/v10"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

var (
	basicRes         context.BasicRes
	apiResources     = map[string]map[string]plugin.ApiResourceHandler{}
	vld              *validator.Validate
	connectionHelper *helper.ConnectionApiHelper
	dsHelper         *helper.DsHelper[models.CopilotConnection, models.CopilotScope, srvhelper.NoScopeConfig]
)

// Init stores basic resources and configures shared helpers for API handlers.
func Init(br context.BasicRes, meta plugin.PluginMeta) {
	basicRes = br
	vld = validator.New()
	connectionHelper = helper.NewConnectionHelper(basicRes, vld, meta.Name())
	dsHelper = helper.NewDataSourceHelper[
		models.CopilotConnection, models.CopilotScope, srvhelper.NoScopeConfig,
	](
		basicRes,
		meta.Name(),
		[]string{"id", "organization"},
		func(c models.CopilotConnection) models.CopilotConnection {
			c.Normalize()
			return c.Sanitize()
		},
		func(s models.CopilotScope) models.CopilotScope { return s },
		nil,
	)

	apiResources = map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": TestConnection,
		},
		"connections": {
			"POST": PostConnections,
			"GET":  ListConnections,
		},
		"connections/:connectionId": {
			"GET":    GetConnection,
			"PATCH":  PatchConnection,
			"DELETE": DeleteConnection,
		},
		"connections/:connectionId/test": {
			"POST": TestExistingConnection,
		},
		"connections/:connectionId/scopes": {
			"GET": GetScopeList,
			"PUT": PutScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    GetScope,
			"PATCH":  PatchScope,
			"DELETE": DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": GetScopeLatestSyncState,
		},
	}
}

// GetApiResources returns registered API handlers for the Copilot plugin.
func GetApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return apiResources
}
