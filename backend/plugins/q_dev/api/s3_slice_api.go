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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
)

type PutScopesReqBody = helper.PutScopesReqBody[models.QDevS3Slice]
type ScopeDetail = srvhelper.ScopeDetail[models.QDevS3Slice, srvhelper.NoScopeConfig]

// PutScopes create or update Q Developer scopes (S3 prefixes)
// @Summary create or update Q Developer scopes
// @Description Create or update Q Developer scopes
// @Tags plugins/q_dev
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.QDevS3Slice
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// GetScopeList returns Q Developer scopes
// @Summary get Q Developer scopes
// @Description get Q Developer scopes
// @Tags plugins/q_dev
// @Param connectionId path int true "connection ID"
// @Param pageSize query int false "page size"
// @Param page query int false "page number"
// @Param blueprints query bool false "include blueprint references"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope returns a single scope record
// @Summary get a Q Developer scope
// @Description get a Q Developer scope
// @Tags plugins/q_dev
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope id"
// @Param blueprints query bool false "include blueprint references"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// PatchScope updates a scope record
// @Summary patch a Q Developer scope
// @Description patch a Q Developer scope
// @Tags plugins/q_dev
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope id"
// @Param scope body models.QDevS3Slice true "json"
// @Success 200  {object} models.QDevS3Slice
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// DeleteScope removes a scope and optionally associated data.
// @Summary delete a Q Developer scope
// @Description delete Q Developer scope data
// @Tags plugins/q_dev
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope id"
// @Param delete_data_only query bool false "Only delete scope data"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}

// GetScopeLatestSyncState returns scope sync state info
// @Summary latest sync state for a Q Developer scope
// @Description get latest sync state for a Q Developer scope
// @Tags plugins/q_dev
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "scope id"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/q_dev/connections/{connectionId}/scopes/{scopeId}/latest-sync-state [GET]
func GetScopeLatestSyncState(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeLatestSyncState(input)
}
