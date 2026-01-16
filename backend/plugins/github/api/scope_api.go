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
	"context"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

type PutScopesReqBody api.PutScopesReqBody[models.GithubRepo]
type ScopeDetail api.ScopeDetail[models.GithubRepo, models.GithubScopeConfig]

// PutScope create or update github repo
// @Summary create or update github repo
// @Description Create or update github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// Enrich missing URL fields by fetching from GitHub API
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	data, ok := input.Body["data"].([]interface{})
	if ok {
		for _, row := range data {
			dict, ok := row.(map[string]interface{})
			if !ok {
				continue
			}

			// Check if URL fields are missing
			htmlUrl, hasHtmlUrl := dict["htmlUrl"].(string)
			cloneUrl, hasCloneUrl := dict["cloneUrl"].(string)
			fullName, hasFullName := dict["fullName"].(string)

			// If URLs are missing but we have fullName, fetch from GitHub API
			if hasFullName && fullName != "" && ((!hasHtmlUrl || htmlUrl == "") || (!hasCloneUrl || cloneUrl == "")) {
				repo, err := getApiRepo(&tasks.GithubOptions{Name: fullName}, apiClient)
				if err != nil {
					// Log warning but don't fail - might be a deleted/private repo
					basicRes.GetLogger().Warn(err, "failed to fetch repo details for %s, URL fields may be incomplete", fullName)
					continue
				}

				// Populate missing fields from API
				if !hasHtmlUrl || htmlUrl == "" {
					dict["htmlUrl"] = repo.HTMLUrl
				}
				if !hasCloneUrl || cloneUrl == "" {
					dict["cloneUrl"] = repo.CloneUrl
				}
			}
		}
	}

	return dsHelper.ScopeApi.PutMultiple(input)
}

// UpdateScope patch to github repo
// @Summary patch to github repo
// @Description patch to github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.GithubRepo true "json"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopeList get Github repos
// @Summary get Github repos
// @Description get Github repos
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes [GET]
func GetScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one Github repo
// @Summary get one Github repo
// @Description get one Github repo
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
