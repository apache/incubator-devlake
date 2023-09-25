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
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// PostScopeConfig create scope config for Github
// @Summary create scope config for Github
// @Description create scope config for Github
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body models.GithubScopeConfig true "scope config"
// @Success 200  {object} models.GithubScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scope-configs [POST]
func PostScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scApi.Post(input)
}

// PatchScopeConfig update scope config for Github
// @Summary update scope config for Github
// @Description update scope config for Github
// @Tags plugins/github
// @Accept application/json
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body models.GithubScopeConfig true "scope config"
// @Success 200  {object} models.GithubScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scope-configs/{id} [PATCH]
func PatchScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	println("api hello")
	fmt.Printf("api patch scope confg %v\n", input.Params)
	return scApi.Patch(input)
}

// GetScopeConfig return one scope config
// @Summary return one scope config
// @Description return one scope config
// @Tags plugins/github
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.GithubScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scope-configs/{id} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scApi.GetDetail(input)
}

// GetScopeConfigList return all scope configs
// @Summary return all scope configs
// @Description return all scope configs
// @Tags plugins/github
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} []models.GithubScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scApi.GetAll(input)
}

// DeleteScopeConfig delete a scope config
// @Summary delete a scope config
// @Description delete a scope config
// @Tags plugins/github
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scope-configs/{id} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scApi.Delete(input)
}
