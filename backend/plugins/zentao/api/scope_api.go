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
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

type PutScopesReqBody api.PutScopesReqBody[models.ZentaoProject]
type ScopeDetail api.ScopeDetail[models.ZentaoProject, models.ZentaoScopeConfig]

// PutScopes create or update zentao projects
// @Summary create or update zentao projects
// @Description Create or update zentao projects
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// PatchScope patch to zentao project
// @Summary patch to zentao project
// @Description patch to zentao project
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.ZentaoProject true "json"
// @Success 200  {object} models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopes get zentao scopes
// @Summary get zentao scopes
// @Description get zentao scopes
// @Tags plugins/zentao
// @Param connectionId path int false "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes [GET]
func GetScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one project
// @Summary get one project
// @Description get one project
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteProjectScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} api.ScopeRefDoc "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteProjectScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
