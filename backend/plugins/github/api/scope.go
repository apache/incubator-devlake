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
	"net/http"
)

// PutScope create or update github repo
// @Summary create or update github repo
// @Description Create or update github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	repos, err := scopeHelper.Put(input)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GithubRepo")
	}
	return &plugin.ApiResourceOutput{Body: repos, Status: http.StatusOK}, nil
}

// UpdateScope patch to github repo
// @Summary patch to github repo
// @Description patch to github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param repoId path int true "repo ID"
// @Param scope body models.GithubRepo true "json"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	repo, err := scopeHelper.Update(input, "github_id")
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, err
	}
	return &plugin.ApiResourceOutput{Body: repo, Status: http.StatusOK}, nil
}

// GetScopeList get Github repos
// @Summary get Github repos
// @Description get Github repos
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []apiRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	apiScopes, err := scopeHelper.GetScopeList(input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: apiScopes, Status: http.StatusOK}, nil
}

// GetScope get one Github repo
// @Summary get one Github repo
// @Description get one Github repo
// @Tags plugins/github
// @Param connectionId path int true "connection ID"
// @Param repoId path int true "repo ID"
// @Success 200  {object} apiRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{id} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	apiScope, err := scopeHelper.GetScope(input, "github_id")
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: apiScope, Status: http.StatusOK}, nil
}
