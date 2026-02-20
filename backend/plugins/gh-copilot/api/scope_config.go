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
)

// PostScopeConfig creates a new scope configuration
// @Summary Create scope config
// @Description Create scope config for GitHub Copilot
// @Tags plugins/gh-copilot
// @Accept json
// @Param connectionId path int true "connection ID"
// @Param request body models.GhCopilotScopeConfig true "scope config"
// @Success 200 {object} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/connections/{connectionId}/scope-configs [POST]
func PostScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Post(input)
}

// GetScopeConfigList returns all scope configurations for a connection
// @Summary Get scope configs
// @Description Get all scope configs for a connection
// @Tags plugins/gh-copilot
// @Param connectionId path int true "connection ID"
// @Success 200 {array} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetAll(input)
}

// GetScopeConfig returns a scope configuration by id
// @Summary Get scope config
// @Description Get a scope config by ID
// @Tags plugins/gh-copilot
// @Param connectionId path int true "connection ID"
// @Param scopeConfigId path int true "scope config ID"
// @Success 200 {object} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/connections/{connectionId}/scope-configs/{scopeConfigId} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetDetail(input)
}

// PatchScopeConfig updates a scope configuration
// @Summary Patch scope config
// @Description Update a scope config
// @Tags plugins/gh-copilot
// @Accept json
// @Param connectionId path int true "connection ID"
// @Param scopeConfigId path int true "scope config ID"
// @Param request body models.GhCopilotScopeConfig true "scope config"
// @Success 200 {object} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/connections/{connectionId}/scope-configs/{scopeConfigId} [PATCH]
func PatchScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Patch(input)
}

// DeleteScopeConfig deletes a scope configuration
// @Summary Delete scope config
// @Description Delete a scope config
// @Tags plugins/gh-copilot
// @Param connectionId path int true "connection ID"
// @Param scopeConfigId path int true "scope config ID"
// @Success 200 {object} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/connections/{connectionId}/scope-configs/{scopeConfigId} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Delete(input)
}

// GetProjectsByScopeConfig returns projects related to a scope config
// @Summary Get projects by scope config
// @Description Get projects details related by scope config
// @Tags plugins/gh-copilot
// @Param scopeConfigId path int true "scope config ID"
// @Success 200 {object} models.GhCopilotScopeConfig
// @Failure 400 {object} shared.ApiBody "bad request"
// @Failure 500 {object} shared.ApiBody "internal error"
// @Router /plugins/gh-copilot/scope-config/{scopeConfigId}/projects [GET]
func GetProjectsByScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetProjectsByScopeConfig(input)
}
