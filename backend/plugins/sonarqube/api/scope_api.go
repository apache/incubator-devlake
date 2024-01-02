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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

type PutScopesReqBody api.PutScopesReqBody[models.SonarqubeProject]
type ScopeDetail api.ScopeDetail[models.SonarqubeProject, srvhelper.NoScopeConfig]

// PutScope create or update sonarqube project
// @Summary create or update sonarqube project
// @Description Create or update sonarqube project
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.SonarqubeProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode request body to scope, deal with lastAnalysisDate format
	data, ok := input.Body["data"].([]interface{})
	if !ok {
		return nil, errors.BadInput.New("invalid `data`")
	}
	for _, item := range data {
		dateStr, ok := item.(map[string]interface{})["lastAnalysisDate"].(string)
		if !ok {
			continue
		}
		timeObj, err := common.ConvertStringToTime(dateStr)
		if err != nil {
			panic(err)
		}

		item.(map[string]interface{})["lastAnalysisDate"] = timeObj

	}

	return dsHelper.ScopeApi.PutMultiple(input)
}

// UpdateScope patch to sonarqube project
// @Summary patch to sonarqube project
// @Description patch to sonarqube project
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "project Key"
// @Param scope body models.SonarqubeProject true "json"
// @Success 200  {object} models.SonarqubeProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopeList get Sonarqube projects
// @Summary get Sonarqube projects
// @Description get Sonarqube projects
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one Sonarqube project
// @Summary get one Sonarqube project
// @Description get one Sonarqube project
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "project key"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/sonarqube
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
