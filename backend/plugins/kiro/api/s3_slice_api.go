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
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/helpers/srvhelper"
	"github.com/apache/devlake/plugins/kiro/models"
)

type PutScopesReqBody = helper.PutScopesReqBody[models.KiroS3Slice]
type ScopeDetail = srvhelper.ScopeDetail[models.KiroS3Slice, srvhelper.NoScopeConfig]

// PutScopes creates or updates Kiro collection scopes.
// @Summary create or update kiro scopes
// @Description Create or update kiro scopes, each covering one AWS account for one month
// @Tags plugins/kiro
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.KiroS3Slice
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// GetScopeList returns the scopes for a connection.
// @Summary get kiro scopes
// @Description get kiro scopes
// @Tags plugins/kiro
// @Param connectionId path int true "connection ID"
// @Param pageSize query int false "page size"
// @Param page query int false "page number"
// @Param blueprints query bool false "include blueprint references"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope returns a single scope.
// @Summary get one kiro scope
// @Description get one kiro scope
// @Tags plugins/kiro
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope ID"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// PatchScope updates a scope.
// @Summary patch a kiro scope
// @Description patch a kiro scope
// @Tags plugins/kiro
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope ID"
// @Param scope body models.KiroS3Slice true "json"
// @Success 200  {object} models.KiroS3Slice
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// DeleteScope removes a scope and optionally its collected data.
// @Summary delete a kiro scope
// @Description delete a kiro scope
// @Tags plugins/kiro
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope ID"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}

// GetScopeLatestSyncState reports the most recent sync for a scope.
// @Summary get the latest sync state of a kiro scope
// @Description get the latest sync state of a kiro scope
// @Tags plugins/kiro
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope ID"
// @Success 200  {object} []models.LatestSyncState
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/scopes/{scopeId}/latest-sync-state [GET]
func GetScopeLatestSyncState(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeLatestSyncState(input)
}
