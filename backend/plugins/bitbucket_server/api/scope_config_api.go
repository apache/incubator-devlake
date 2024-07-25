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
)

// CreateScopeConfig create scope config for Bitbucket
// @Summary create scope config for Bitbucket
// @Description create scope config for Bitbucket
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body models.BitbucketServerScopeConfig true "scope config"
// @Success 200  {object} models.BitbucketServerScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/scope-configs [POST]
func CreateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Post(input)
}

// UpdateScopeConfig update scope config for Bitbucket
// @Summary update scope config for Bitbucket
// @Description update scope config for Bitbucket
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param scopeConfigId path int true "scopeConfigId"
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body models.BitbucketServerScopeConfig true "scope config"
// @Success 200  {object} models.BitbucketServerScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/scope-configs/{scopeConfigId} [PATCH]
func UpdateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeConfigId"] = strings.TrimLeft(input.Params["scopeConfigId"], "/")
	return dsHelper.ScopeConfigApi.Patch(input)
}

// GetScopeConfig return one scope config
// @Summary return one scope config
// @Description return one scope config
// @Tags plugins/bitbucket_server
// @Param scopeConfigId path int true "scopeConfigId"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.BitbucketServerScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/scope-configs/{scopeConfigId} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeConfigId"] = strings.TrimLeft(input.Params["scopeConfigId"], "/")
	return dsHelper.ScopeConfigApi.GetDetail(input)
}

// GetScopeConfigList return all scope configs
// @Summary return all scope configs
// @Description return all scope configs
// @Tags plugins/bitbucket_server
// @Param connectionId path int true "connectionId"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []models.BitbucketServerScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetAll(input)
}

// GetProjectsByScopeConfig return projects details related by scope config
// @Summary return all related projects
// @Description return all related projects
// @Tags plugins/bitbucket_server
// @Param id path int true "id"
// @Param scopeConfigId path int true "scopeConfigId"
// @Success 200  {object} models.ProjectScopeOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/scope-config/{scopeConfigId}/projects [GET]
func GetProjectsByScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetProjectsByScopeConfig(input)
}

// DeleteScopeConfig delete a scope config
// @Summary delete a scope config
// @Description delete a scope config
// @Tags plugins/bitbucket_server
// @Param scopeConfigId path int true "scopeConfigId"
// @Param connectionId path int true "connectionId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/scope-configs/{scopeConfigId} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeConfigId"] = strings.TrimLeft(input.Params["scopeConfigId"], "/")
	return dsHelper.ScopeConfigApi.Delete(input)
}
