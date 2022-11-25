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
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

// PutScope create or update github repo
// @Summary create or update github repo
// @Description Create or update github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param repoId path int false "repo ID"
// @Param scope body models.GithubRepo true "json"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{repoId} [PUT]
func PutScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, repoId := extractParam(input.Params)
	if connectionId*repoId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or repoId")
	}
	var repo models.GithubRepo
	err := errors.Convert(mapstructure.Decode(input.Body, &repo))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Github repo error")
	}
	err = verifyRepo(&repo)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().CreateOrUpdate(repo)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GithubRepo")
	}
	return &core.ApiResourceOutput{Body: repo, Status: http.StatusOK}, nil
}

// UpdateScope patch to github repo
// @Summary patch to github repo
// @Description patch to github repo
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param repoId path int false "repo ID"
// @Param scope body models.GithubRepo true "json"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{repoId} [PATCH]
func UpdateScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, repoId := extractParam(input.Params)
	if connectionId*repoId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or repoId")
	}
	var repo models.GithubRepo
	err := basicRes.GetDal().First(&repo, dal.Where("connection_id = ? AND github_id = ?", connectionId, repoId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting GithubRepo error")
	}
	err = helper.DecodeMapStruct(input.Body, &repo)
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
	return &core.ApiResourceOutput{Body: repo, Status: http.StatusOK}, nil
}

// GetScopeList get Github repos
// @Summary get Github repos
// @Description get Github repos
// @Tags plugins/github
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var repos []models.GithubRepo
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := basicRes.GetDal().All(&repos, dal.Where("connection_id = ?", connectionId))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: repos, Status: http.StatusOK}, nil
}

// GetScope get one Github repo
// @Summary get one Github repo
// @Description get one Github repo
// @Tags plugins/github
// @Param connectionId path int false "connection ID"
// @Param repoId path int false "repo ID"
// @Success 200  {object} models.GithubRepo
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/scopes/{repoId} [GET]
func GetScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var repo models.GithubRepo
	connectionId, repoId := extractParam(input.Params)
	if connectionId*repoId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := basicRes.GetDal().First(&repo, dal.Where("connection_id = ? AND github_id = ?", connectionId, repoId))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: repo, Status: http.StatusOK}, nil
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
	if repo.GithubId == 0 {
		return errors.BadInput.New("invalid repoId")
	}
	if repo.ScopeId != strconv.Itoa(repo.GithubId) {
		return errors.BadInput.New("the scope_id does not match the github_id")
	}
	return nil
}
