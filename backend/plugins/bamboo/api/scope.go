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
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

type scopeReq struct {
	Data []models.BambooPlan `json:"data"`
}

// PutScope create or update bamboo plan
// @Summary create or update bamboo plan
// @Description Create or update bamboo plan
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body scopeReq true "json"
// @Success 200  {object} models.BambooPlan
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Put(input)
}

// UpdateScope patch to bamboo plan
// @Summary patch to bamboo plan
// @Description patch to bamboo plan
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "plan key"
// @Param scope body models.BambooPlan true "json"
// @Success 200  {object} models.BambooPlan
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input)
}

// GetScopeList get Bamboo plans
// @Summary get Bamboo plans
// @Description get Bamboo plans
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []models.BambooPlan
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

// GetScope get one Bamboo plan
// @Summary get one Bamboo plan
// @Description get one Bamboo plan
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "plan key"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} models.BambooPlan
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScope(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/bamboo
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} api.ScopeRefDoc "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Delete(input)
}
