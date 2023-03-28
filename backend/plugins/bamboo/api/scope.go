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
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

type ScopeRes struct {
	models.BambooProject
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type ScopeReq api.ScopeReq[models.BambooProject]

// PutScope create or update bamboo project
// @Summary create or update bamboo project
// @Description Create or update bamboo project
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body ScopeReq true "json"
// @Success 200  {object} []models.BambooProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Put(input)
}

// UpdateScope patch to bamboo project
// @Summary patch to bamboo project
// @Description patch to bamboo project
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "project ID"
// @Param scope body models.BambooProject true "json"
// @Success 200  {object} models.BambooProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input, "project_key")
}

// GetScopeList get Bamboo projects
// @Summary get Bamboo projects
// @Description get Bamboo projects
// @Tags plugins/bamboo
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

// GetScope get one Bamboo project
// @Summary get one Bamboo project
// @Description get one Bamboo project
// @Tags plugins/bamboo
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "project ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScope(input, "project_key")
}
