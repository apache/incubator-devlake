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
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func listGithubRemoteScopes(
	connection *models.GithubConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page GithubRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo],
	nextPage *GithubRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	if page.PerPage == 0 {
		page.PerPage = 100
	}

	if connection.AuthMethod == plugin.AUTH_METHOD_APPKEY {
		return listGithubAppInstalledRepos(apiClient, page)
	}
	if groupId == "" {
		return listGithubUserOrgs(apiClient, page)
	}
	return listGithubOrgRepos(apiClient, groupId, page)
}

func listGithubUserOrgs(
	apiClient plugin.ApiClient,
	page GithubRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo],
	nextPage *GithubRemotePagination,
	err errors.Error,
) {
	// user's own org
	if page.Page == 1 {
		userBody, err := apiClient.Get("user", nil, nil)
		if err != nil {
			return nil, nil, err
		}
		var o org
		errors.Must(api.UnmarshalResponse(userBody, &o))
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       fmt.Sprintf("%v", o.Login),
			Name:     fmt.Sprintf("%v", o.Login),
			FullName: fmt.Sprintf("%v", o.Login),
		})
	}
	// user's orgs
	orgsBody, err := apiClient.Get(
		"user/orgs",
		url.Values{
			"page":     []string{fmt.Sprintf("%v", page.Page)},
			"per_page": []string{fmt.Sprintf("%v", page.PerPage)},
		},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	if orgsBody.StatusCode != http.StatusOK {
		err := errors.Convert(tryToResolveErrorMessage(orgsBody))
		return nil, nil, err
	}
	var orgs []org
	if err := api.UnmarshalResponse(orgsBody, &orgs); err != nil {
		return nil, nil, err
	}
	for _, o := range orgs {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       fmt.Sprintf("%v", o.Login),
			Name:     fmt.Sprintf("%v", o.Login),
			FullName: fmt.Sprintf("%v", o.Login),
		})
	}
	// there may be more orgs
	if len(orgs) == page.PerPage {
		nextPage = &GithubRemotePagination{
			Page:    page.Page + 1,
			PerPage: page.PerPage,
		}
	}
	return children, nextPage, nil
}

func tryToResolveErrorMessage(resp *http.Response) error {
	if resp == nil {
		return goerrors.New("nil response")
	}
	type Response struct {
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url"`
	}
	resBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if len(resBody) == 0 {
		return goerrors.New("empty response")
	}
	var respStruct Response
	if err := json.Unmarshal(resBody, &respStruct); err != nil {
		return err
	}
	return fmt.Errorf("response status code: %d, message: %s", resp.StatusCode, respStruct.Message)
}

func listGithubOrgRepos(
	apiClient plugin.ApiClient,
	org string,
	page GithubRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo],
	nextPage *GithubRemotePagination,
	err errors.Error,
) {
	query := url.Values{
		"page":     []string{fmt.Sprintf("%v", page.Page)},
		"per_page": []string{fmt.Sprintf("%v", page.PerPage)},
	}
	// user's orgs
	reposBody, err := apiClient.Get(fmt.Sprintf("orgs/%s/repos", org), query, nil)
	if err != nil {
		return nil, nil, err
	}
	// if not found, try to get user repos
	if reposBody.StatusCode == http.StatusNotFound {
		reposBody, err = apiClient.Get(fmt.Sprintf("users/%s/repos", org), query, nil)
		if err != nil {
			return nil, nil, err
		}
	}
	var repos []repo
	errors.Must(api.UnmarshalResponse(reposBody, &repos))
	for _, repo := range repos {
		children = append(children, toRepoModel(&repo))
	}
	// there may be more repos
	if len(children) == page.PerPage {
		nextPage = &GithubRemotePagination{
			Page:    page.Page + 1,
			PerPage: page.PerPage,
		}
	}
	return children, nextPage, nil
}

func listGithubAppInstalledRepos(
	apiClient plugin.ApiClient,
	page GithubRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo],
	nextPage *GithubRemotePagination,
	err errors.Error,
) {
	resApp, err := apiClient.Get("installation/repositories",
		url.Values{
			"page":     []string{fmt.Sprintf("%v", page.Page)},
			"per_page": []string{fmt.Sprintf("%v", page.PerPage)},
		},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	var appRepos GithubAppRepo
	errors.Must(api.UnmarshalResponse(resApp, &appRepos))
	for _, r := range appRepos.Repositories {
		orgName := r.Owner.Login
		if orgName != "" {
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo]{
				Type:     api.RAS_ENTRY_TYPE_SCOPE,
				ParentId: &orgName,
				Id:       fmt.Sprintf("%v", r.ID),
				Name:     fmt.Sprintf("%v", r.Name),
				FullName: fmt.Sprintf("%v", r.FullName),
				Data: &models.GithubRepo{
					GithubId:    r.ID,
					Name:        r.Name,
					FullName:    r.FullName,
					HTMLUrl:     r.HTMLURL,
					Description: r.Description,
					OwnerId:     r.Owner.ID,
					CloneUrl:    r.CloneURL,
				},
			})
		}
	}
	if len(appRepos.Repositories) == page.PerPage {
		nextPage = &GithubRemotePagination{
			Page:    page.Page + 1,
			PerPage: page.PerPage,
		}
	}
	return
}

func searchGithubRepos(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo],
	err errors.Error,
) {
	res, err := apiClient.Get(
		"search/repositories",
		url.Values{
			"q":        []string{params.Search},
			"page":     []string{fmt.Sprintf("%v", params.Page)},
			"per_page": []string{fmt.Sprintf("%v", params.PageSize)},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var resBody struct {
		Items []repo `json:"items"`
	}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	for _, r := range resBody.Items {
		children = append(children, toRepoModel(&r))
	}
	return
}

func toRepoModel(r *repo) dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo] {
	parentId := fmt.Sprintf("%v", r.Owner.Login)
	return dsmodels.DsRemoteApiScopeListEntry[models.GithubRepo]{
		Type:     api.RAS_ENTRY_TYPE_SCOPE,
		ParentId: &parentId,
		Id:       fmt.Sprintf("%v", r.ID),
		Name:     fmt.Sprintf("%v", r.Name),
		FullName: fmt.Sprintf("%v", r.FullName),
		Data: &models.GithubRepo{
			GithubId:    r.ID,
			Name:        r.Name,
			FullName:    r.FullName,
			HTMLUrl:     r.HTMLURL,
			Description: r.Description,
			OwnerId:     r.Owner.ID,
			CloneUrl:    r.CloneURL,
			CreatedDate: r.CreatedAt,
			UpdatedDate: r.UpdatedAt,
		},
	}
}

// RemoteScopes list all available scopes on the remote server
// @Summary list all available scopes on the remote server
// @Description list all available scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GithubRepo]
// @Tags plugins/github
// @Router /plugins/github/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes searches scopes on the remote server
// @Summary searches scopes on the remote server
// @Description searches scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.GithubRepo] "the parentIds are always null"
// @Tags plugins/github
// @Router /plugins/github/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/github
// @Router /plugins/github/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
