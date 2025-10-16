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

// CreateScopeConfig create scope config for ArgoCD
// @Summary create scope config for ArgoCD
// @Description create scope config for ArgoCD
// @Tags plugins/argocd
// @Accept application/json
// @Param connectionId path int false "connectionId"
// @Param scopeConfig body models.ArgocdScopeConfig true "scope config"
// @Success 200  {object} models.ArgocdScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scope-configs [POST]
func CreateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Post(input)
}

// UpdateScopeConfig update scope config for ArgoCD
// @Summary update scope config for ArgoCD
// @Description update scope config for ArgoCD
// @Tags plugins/argocd
// @Accept application/json
// @Param id path int true "id"
// @Param connectionId path int false "connectionId"
// @Param scopeConfig body models.ArgocdScopeConfig true "scope config"
// @Success 200  {object} models.ArgocdScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scope-configs/{id} [PATCH]
func UpdateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Patch(input)
}

// GetScopeConfigList get scope config list for ArgoCD
// @Summary get scope config list for ArgoCD
// @Description get scope config list for ArgoCD
// @Tags plugins/argocd
// @Param connectionId path int false "connectionId"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page number, default 1"
// @Success 200  {object} []models.ArgocdScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetAll(input)
}

// GetScopeConfig get scope config for ArgoCD
// @Summary get scope config for ArgoCD
// @Description get scope config for ArgoCD
// @Tags plugins/argocd
// @Param id path int true "id"
// @Param connectionId path int false "connectionId"
// @Success 200  {object} models.ArgocdScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scope-configs/{id} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetDetail(input)
}

// DeleteScopeConfig delete scope config for ArgoCD
// @Summary delete scope config for ArgoCD
// @Description delete scope config for ArgoCD
// @Tags plugins/argocd
// @Param id path int true "id"
// @Param connectionId path int false "connectionId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scope-configs/{id} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.Delete(input)
}

// GetProjectsByScopeConfig get projects related to a scope config
// @Summary get projects related to scope config
// @Description get projects related to scope config
// @Tags plugins/argocd
// @Param scopeConfigId path int true "scopeConfigId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/scope-config/{scopeConfigId}/projects [GET]
func GetProjectsByScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeConfigApi.GetProjectsByScopeConfig(input)
}
