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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
)

type PutScopesReqBody api.PutScopesReqBody[models.TeambitionProject]
type ScopeDetail api.ScopeDetail[models.TeambitionProject, models.TeambitionScopeConfig]

// PutScopes create or update Azure DevOps repo
// @Summary create or update Azure DevOps repo
// @Description Create or update Azure DevOps repo
// @Tags plugins/teambition
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.TeambitionProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// PatchScope patch to Azure DevOps repo
// @Summary patch to Azure DevOps repo
// @Description patch to Azure DevOps repo
// @Tags plugins/teambition
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.TeambitionProject true "json"
// @Success 200  {object} models.TeambitionProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopes get Azure DevOps repos
// @Summary get Azure DevOps repos
// @Description get Azure DevOps repos
// @Tags plugins/teambition
// @Param connectionId path int true "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/scopes [GET]
func GetScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one Azure DevOps repo
// @Summary get one Azure DevOps repo
// @Description get one Azure DevOps repo
// @Tags plugins/teambition
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/teambition
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200  {object} models.TeambitionProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
