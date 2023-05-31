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

// CreateScopeConfig create scope config for Trello
// @Summary create scope config for Trello
// @Description create scope config for Trello
// @Tags plugins/trello
// @Accept application/json
// @Param scopeConfig body models.TrelloScopeConfig true "scope config"
// @Success 200  {object} models.TrelloScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scope_configs [POST]
func CreateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.Create(input)
}

// UpdateScopeConfig update scope config for Trello
// @Summary update scope config for Trello
// @Description update scope config for Trello
// @Tags plugins/trello
// @Accept application/json
// @Param id path int true "id"
// @Param scopeConfig body models.TrelloScopeConfig true "scope config"
// @Success 200  {object} models.TrelloScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scope_configs/{id} [PATCH]
func UpdateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.Update(input)
}

// GetScopeConfig return one scope config
// @Summary return one scope config
// @Description return one scope config
// @Tags plugins/trello
// @Param id path int true "id"
// @Success 200  {object} models.TrelloScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scope_configs/{id} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.Get(input)
}

// GetScopeConfigList return all scope configs
// @Summary return all scope configs
// @Description return all scope configs
// @Tags plugins/trello
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []models.TrelloScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scope_configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.List(input)
}
