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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
)

type apiRepo struct {
	models.GithubRepo
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type req struct {
	Data []*models.GithubRepo `json:"data"`
}

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
	var repos req
	err := errors.Convert(mapstructure.Decode(input.Body, &repos))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Github repo error")
	}
	fmt.Println("length is ", len(repos.Data))
	err = scopeHelper.Put(input, "githubId", repos.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GithubRepo")
	}
	return &plugin.ApiResourceOutput{Body: repos.Data, Status: http.StatusOK}, nil
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
// @Router /plugins/github/connections/{connectionId}/scopes/{repoId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, repoId := extractParam(input.Params)
	if connectionId*repoId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or repoId")
	}
	var repo models.GithubRepo
	err := basicRes.GetDal().First(&repo, dal.Where("connection_id = ? AND github_id = ?", connectionId, repoId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting GithubRepo error")
	}
	err = api.DecodeMapStruct(input.Body, &repo)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch github repo error")
	}
	err = verifyRepo(&repo)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().Update(repo)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GithubRepo")
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
	var repos []models.GithubRepo
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&repos, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var ruleIds []uint64
	for _, repo := range repos {
		if repo.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, repo.TransformationRuleId)
		}
	}
	var rules []models.GithubTransformationRule
	if len(ruleIds) > 0 {
		err = basicRes.GetDal().All(&rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		names[rule.ID] = rule.Name
	}
	var apiRepos []apiRepo
	for _, repo := range repos {
		apiRepos = append(apiRepos, apiRepo{repo, names[repo.TransformationRuleId]})
	}
	return &plugin.ApiResourceOutput{Body: apiRepos, Status: http.StatusOK}, nil
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
// @Router /plugins/github/connections/{connectionId}/scopes/{repoId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var repo models.GithubRepo
	connectionId, repoId := extractParam(input.Params)
	if connectionId*repoId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	db := basicRes.GetDal()
	err := db.First(&repo, dal.Where("connection_id = ? AND github_id = ?", connectionId, repoId))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	var rule models.GithubTransformationRule
	if repo.TransformationRuleId > 0 {
		err = basicRes.GetDal().First(&rule, dal.Where("id = ?", repo.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}
	return &plugin.ApiResourceOutput{Body: apiRepo{repo, rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	repoId, _ := strconv.ParseUint(params["repoId"], 10, 64)
	return connectionId, repoId
}

func verifyRepo(repo *models.GithubRepo) errors.Error {
	if repo.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if repo.GithubId <= 0 {
		return errors.BadInput.New("invalid github ID")
	}
	return nil
}
