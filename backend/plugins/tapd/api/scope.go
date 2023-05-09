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
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"strings"
)

type ScopeRes struct {
	models.TapdWorkspace
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type TapdScopeReq api.ScopeReq[models.TapdWorkspace]

// PutScope create or update tapd job
// @Summary create or update tapd job
// @Description Create or update tapd job
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body TapdScopeReq true "json"
// @Success 200  {object} []models.TapdWorkspace
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Put(input)
}

// UpdateScope patch to tapd job
// @Summary patch to tapd job
// @Description patch to tapd job
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "workspace ID"
// @Param scope body models.TapdWorkspace true "json"
// @Success 200  {object} models.TapdWorkspace
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input, "id")
}

// GetScopeList get tapd jobs
// @Summary get tapd jobs
// @Description get tapd jobs
// @Tags plugins/tapd
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

// GetScope get one tapd job
// @Summary get one tapd job
// @Description get one tapd job
// @Tags plugins/tapd
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "workspace ID"
// @Success 200  {object} ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	input.Params["scopeId"] = strings.TrimLeft(input.Params["scopeId"], "/")
	return scopeHelper.GetScope(input, "id")
}
