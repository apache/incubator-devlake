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
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/mitchellh/mapstructure"
)

// CreateScopeConfig create scope config for Jira
// @Summary create scope config for Jira
// @Description create scope config for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body tasks.JiraScopeConfig true "scope config"
// @Success 200  {object} tasks.JiraScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scope-configs [POST]
func CreateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	rule, err := makeDbScopeConfigFromInput(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in makeJiraScopeConfig")
	}
	newRule := map[string]interface{}{}
	err = errors.Convert(mapstructure.Decode(rule, &newRule))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in makeJiraScopeConfig")
	}
	input.Body = newRule
	return scHelper.Create(input)
}

// UpdateScopeConfig update scope config for Jira
// @Summary update scope config for Jira
// @Description update scope config for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body tasks.JiraScopeConfig true "scope config"
// @Success 200  {object} tasks.JiraScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scope-configs/{id} [PATCH]
func UpdateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	scopeConfigId, e := strconv.ParseUint(input.Params["id"], 10, 64)
	if e != nil {
		return nil, errors.Default.Wrap(e, "the scope config ID should be an integer")
	}
	var req tasks.JiraScopeConfig
	err := api.Decode(input.Body, &req, vld)
	if err != nil {
		return nil, err
	}
	var oldDB models.JiraScopeConfig
	err = basicRes.GetDal().First(&oldDB, dal.Where("id = ?", scopeConfigId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on getting ScopeConfig")
	}
	oldTr, err := tasks.MakeScopeConfig(oldDB)
	if err != nil {
		return nil, err
	}
	err = api.DecodeMapStruct(input.Body, oldTr, true)
	if err != nil {
		return nil, err
	}

	newDB, err := oldTr.ToDb()
	if err != nil {
		return nil, err
	}
	newDB.ID = scopeConfigId
	newDB.ConnectionId = connectionId
	newDB.CreatedAt = oldDB.CreatedAt
	err = basicRes.GetDal().Update(newDB)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: newDB, Status: http.StatusOK}, err
}

func makeDbScopeConfigFromInput(input *plugin.ApiResourceInput) (*models.JiraScopeConfig, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var req tasks.JiraScopeConfig
	err := api.Decode(input.Body, &req, vld)
	if err != nil {
		return nil, err
	}
	req.ConnectionId = connectionId
	return req.ToDb()
}

// GetScopeConfig return one scope config
// @Summary return one scope config
// @Description return one scope config
// @Tags plugins/jira
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} tasks.JiraScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scope-configs/{id} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.Get(input)
}

// GetScopeConfigList return all scope configs
// @Summary return all scope configs
// @Description return all scope configs
// @Tags plugins/jira
// @Param connectionId path int true "connectionId"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []tasks.JiraScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.List(input)
}

// GetApplicationTypes return issue application types
// @Summary return issue application types
// @Description return issue application types
// @Tags plugins/jira
// @Param connectionId path int true "connectionId"
// @Param key query string false "issue key"
// @Success 200  {object} []string
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/application-types [GET]
func GetApplicationTypes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.JiraConnection
	err := connectionHelper.First(&connection, input.Params)
	if err != nil {
		return nil, err
	}
	key := input.Query.Get("key")
	if key == "" {
		return nil, errors.BadInput.New("key is empty")
	}

	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	res, err = apiClient.Get(fmt.Sprintf("api/2/issue/%s", key), nil, nil)
	if err != nil {
		return nil, err
	}
	var issue struct {
		Id string `json:"id"`
	}
	err = api.UnmarshalResponse(res, &issue)
	if err != nil {
		return nil, err
	}
	// get application types
	query := url.Values{}
	query.Set("issueId", issue.Id)
	res, err = apiClient.Get("dev-status/1.0/issue/summary", query, nil)
	if err != nil {
		return nil, err
	}
	var summary struct {
		Summary struct {
			Repository struct {
				ByInstanceType map[string]interface{} `json:"byInstanceType"`
			} `json:"repository"`
		} `json:"summary"`
	}
	err = api.UnmarshalResponse(res, &summary)
	if err != nil {
		return nil, err
	}
	var types []string
	for k := range summary.Summary.Repository.ByInstanceType {
		types = append(types, k)
	}
	sort.Strings(types)
	return &plugin.ApiResourceOutput{Body: types, Status: http.StatusOK}, nil
}

// GetCommitsURLs return some commits URLs
// @Summary return some commits URLs, at most 5
// @Description return some commits URLs, at most 5
// @Tags plugins/jira
// @Param connectionId path int true "connectionId"
// @Param key query string true "issue key"
// @Param applicationType query string true "application type"
// @Success 200  {object} []string
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/dev-panel-commits [GET]
func GetCommitsURLs(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.JiraConnection
	err := connectionHelper.First(&connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get issue key
	key := input.Query.Get("key")
	if key == "" {
		return nil, errors.BadInput.New("key is empty")
	}
	// get application types
	applicationType := input.Query.Get("applicationType")
	if applicationType == "" {
		return nil, errors.BadInput.New("applicationType is empty")
	}
	// get issue ID from issue key
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	res, err = apiClient.Get(fmt.Sprintf("api/2/issue/%s", key), nil, nil)
	if err != nil {
		return nil, err
	}
	var issue struct {
		Id string `json:"id"`
	}
	err = api.UnmarshalResponse(res, &issue)
	if err != nil {
		return nil, err
	}
	// get commits
	query := url.Values{}
	query.Set("issueId", issue.Id)
	query.Set("applicationType", applicationType)
	query.Set("dataType", "repository")
	res, err = apiClient.Get("dev-status/1.0/issue/detail", query, nil)
	if err != nil {
		return nil, err
	}
	type commit struct {
		ID              string          `json:"id"`
		DisplayID       string          `json:"displayId"`
		AuthorTimestamp api.Iso8601Time `json:"authorTimestamp"`
		URL             string          `json:"url"`
	}
	var detail struct {
		Detail []struct {
			Repositories []struct {
				Commits []commit `json:"commits"`
			} `json:"repositories"`
		} `json:"detail"`
	}

	err = api.UnmarshalResponse(res, &detail)
	if err != nil {
		return nil, err
	}
	var commits []commit
	for _, item := range detail.Detail {
		for _, repo := range item.Repositories {
			commits = append(commits, repo.Commits...)
		}
	}

	// sort by authorTimestamp
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].AuthorTimestamp.ToTime().After(commits[j].AuthorTimestamp.ToTime())
	})
	// return at most 5 commits
	var commitURLs []string
	for i, cmt := range commits {
		if i >= 5 {
			break
		}
		commitURLs = append(commitURLs, cmt.URL)
	}
	return &plugin.ApiResourceOutput{Body: commitURLs, Status: http.StatusOK}, nil
}
