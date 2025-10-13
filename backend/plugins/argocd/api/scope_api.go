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
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// PutScopes create or update argocd applications
// @Summary create or update argocd applications
// @Description Create or update ArgoCD application scopes
// @Tags plugins/argocd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body models.ArgocdApplication true "json"
// @Success 200  {object} []models.ArgocdApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connParam, ok := input.Params["connectionId"]
	if ok {
		if cid, err := strconv.ParseUint(connParam, 10, 64); err == nil {
			cfg, _ := CreateDefaultScopeConfig(cid)
			if cfg != nil {
				if data, ok := input.Body["data"].([]interface{}); ok {
					for _, row := range data {
						if m, ok := row.(map[string]interface{}); ok {
							if _, has := m["scopeConfigId"]; !has || m["scopeConfigId"] == 0 {
								m["scopeConfigId"] = cfg.ID
							}
						}
					}
				}
			}
		}
	}
	return dsHelper.ScopeApi.PutMultiple(input)
}

// UpdateScope patch to argocd application
// @Summary patch to argocd application
// @Description Patch ArgoCD application scope
// @Tags plugins/argocd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "application name"
// @Param scope body models.ArgocdApplication true "json"
// @Success 200  {object} models.ArgocdApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopeList get ArgoCD applications
// @Summary get ArgoCD applications
// @Description Get ArgoCD applications
// @Tags plugins/argocd
// @Param connectionId path int false "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []models.ArgocdApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one ArgoCD application
// @Summary get one ArgoCD application
// @Description Get one ArgoCD application
// @Tags plugins/argocd
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "application name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page number, default 1"
// @Success 200  {object} models.ArgocdApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description Delete data associated with ArgoCD application scope
// @Tags plugins/argocd
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "application name"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
