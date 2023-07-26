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
	"regexp"
	"sort"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type genRegexReq struct {
	Pattern string `json:"pattern"`
}

type genRegexResp struct {
	Regex string `json:"regex"`
}

type applyRegexReq struct {
	Regex string   `json:"regex"`
	Urls  []string `json:"urls"`
}

type repo struct {
	Namespace string `json:"namespace"`
	RepoName  string `json:"repo_name"`
	CommitSha string `json:"commit_sha"`
}

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
	return scHelper.Update(input)
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

// DeleteScopeConfig delete a scope config
// @Summary delete a scope config
// @Description delete a scope config
// @Tags plugins/jira
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scope-configs/{id} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scHelper.Delete(input)
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

// GenRegex generate regex from url
// @Summary generate regex from url
// @Description generate regex from url
// @Tags plugins/jira
// @Param generate-regex body genRegexReq true "generate regex"
// @Success 200  {object} genRegexResp
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/generate-regex [POST]
func GenRegex(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req genRegexReq
	err := api.Decode(input.Body, &req, nil)
	if err != nil {
		return nil, err
	}
	err = checkInput(req.Pattern)
	if err != nil {
		return nil, err
	}
	reg := genRegex(req.Pattern)
	_, e := regexp.Compile(reg)
	if e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid url")
	}

	return &plugin.ApiResourceOutput{Body: genRegexResp{Regex: reg}, Status: http.StatusOK}, nil
}

func checkInput(input string) errors.Error {
	input = strings.TrimSpace(input)
	if input == "" {
		return errors.BadInput.New("empty input")
	}
	if !strings.Contains(input, "{namespace}") {
		return errors.BadInput.New("missing {namespace}")
	}
	if !strings.Contains(input, "{repo_name}") {
		return errors.BadInput.New("missing {repo_name}")
	}
	if !strings.Contains(input, "{commit_sha}") {
		return errors.BadInput.New("missing {commit_sha}")
	}
	return nil
}

func genRegex(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "{namespace}", `(?P<namespace>\S+)`, -1)
	s = strings.Replace(s, "{repo_name}", `(?P<repo_name>\S+)`, -1)
	s = strings.Replace(s, "{commit_sha}", `(?P<commit_sha>\w{40})`, -1)
	return s
}

// ApplyRegex return parsed commits URLs
// @Summary return parsed commits URLs
// @Description return parsed commits URLs
// @Tags plugins/jira
// @Param apply-regex body applyRegexReq true "apply regex"
// @Success 200  {object} []string
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/apply-regex [POST]
func ApplyRegex(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req applyRegexReq
	err := api.Decode(input.Body, &req, nil)
	if err != nil {
		return nil, err
	}
	var repos []*repo
	for _, u := range req.Urls {
		r, e1 := applyRegex(req.Regex, u)
		if e1 != nil {
			return nil, err
		}
		repos = append(repos, r)
	}
	return &plugin.ApiResourceOutput{Body: repos, Status: http.StatusOK}, nil
}

func applyRegex(regexStr, commitUrl string) (*repo, errors.Error) {
	pattern, e := regexp.Compile(regexStr)
	if e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid regex")
	}
	if !pattern.MatchString(commitUrl) {
		return nil, errors.BadInput.New("invalid url")
	}
	group := pattern.FindStringSubmatch(commitUrl)
	if len(group) != 4 {
		return nil, errors.BadInput.New("invalid group count")
	}
	r := new(repo)
	for i, name := range pattern.SubexpNames() {
		if i != 0 && name != "" {
			switch name {
			case "namespace":
				r.Namespace = group[i]
			case "repo_name":
				r.RepoName = group[i]
			case "commit_sha":
				r.CommitSha = group[i]
			default:
				return nil, errors.BadInput.New("invalid group name")
			}
		}
	}
	return r, nil
}
