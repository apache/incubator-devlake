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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

type PutScopesReqBody api.PutScopesReqBody[models.BitbucketRepo]
type ScopeDetail api.ScopeDetail[models.BitbucketRepo, models.BitbucketScopeConfig]

// PutScopes create or update repo
// @Summary create or update repo
// @Description Create or update repo
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.BitbucketRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// PatchScope patch to repo
// @Summary patch to repo
// @Description patch to repo
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "repo ID"
// @Param scope body models.BitbucketRepo true "json"
// @Success 200  {object} models.BitbucketRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeId"] = strings.TrimLeft(input.Params["scopeId"], "/")
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopes get repos
// @Summary get repos
// @Description get repos
// @Tags plugins/bitbucket
// @Param connectionId path int true "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/scopes/ [GET]
func GetScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

func GetScopeDispatcher(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeIdWithSuffix := strings.TrimLeft(input.Params["scopeId"], "/")
	if strings.HasSuffix(scopeIdWithSuffix, "/latest-sync-state") {
		input.Params["scopeId"] = strings.TrimSuffix(scopeIdWithSuffix, "/latest-sync-state")
		return GetScopeLatestSyncState(input)
	}
	return GetScope(input)
}

// GetScope get one repo
// @Summary get one repo
// @Description get one repo
// @Tags plugins/bitbucket
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "repo ID"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeId"] = strings.TrimLeft(input.Params["scopeId"], "/")
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/bitbucket
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} api.ScopeRefDoc "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeId"] = strings.TrimLeft(input.Params["scopeId"], "/")
	return dsHelper.ScopeApi.Delete(input)
}
