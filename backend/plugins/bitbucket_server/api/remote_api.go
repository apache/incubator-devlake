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
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

func listBitbucketServerRemoteScopes(
	connection *models.BitbucketServerConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page BitBucketServerRemotePagination) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo],
	*BitBucketServerRemotePagination,
	errors.Error,
) {
	if page.Limit == 0 {
		page.Limit = 100
	}

	if groupId == "" {
		return listBitbucketServerProjects(apiClient, page)
	}

	return listBitbucketServerRepos(apiClient, groupId, page)
}

func listBitbucketServerProjects(apiClient plugin.ApiClient, page BitBucketServerRemotePagination) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo],
	*BitBucketServerRemotePagination,
	errors.Error,
) {
	children := []dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{}

	for {
		query := initialQuery(page)

		res, err := apiClient.Get("rest/api/1.0/projects", query, nil)
		if err != nil {
			return nil, nil, err
		}

		resBody := &models.ProjectsResponse{}
		err = api.UnmarshalResponse(res, resBody)
		if err != nil {
			return nil, nil, err
		}

		for _, r := range resBody.Values {
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{
				Type:     api.RAS_ENTRY_TYPE_GROUP,
				Id:       fmt.Sprintf("%v", r.Key),
				ParentId: nil,
				Name:     r.Name,
				FullName: r.Name,
			})
		}

		if resBody.IsLastPage {
			break
		}

		if resBody.NextPageStart == nil {
			if len(resBody.Values) >= page.Limit {
				page.Start += len(resBody.Values)
			} else {
				break
			}
		} else {
			if *resBody.NextPageStart < page.Start {
				break
			}
			page.Start = *resBody.NextPageStart
		}
	}

	return children, nil, nil
}

func listBitbucketServerRepos(apiClient plugin.ApiClient, groupId string, page BitBucketServerRemotePagination) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo],
	*BitBucketServerRemotePagination,
	errors.Error,
) {
	if groupId == "" {
		return nil, nil, nil
	}

	query := initialQuery(page)

	res, err := apiClient.Get(fmt.Sprintf("rest/api/1.0/projects/%s/repos", groupId), query, nil)
	if err != nil {
		return nil, nil, err
	}

	resBody := &models.ReposResponse{}
	err = api.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, nil, err
	}

	children := []dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{}
	for _, r := range resBody.Values {
		parent := r.Project.Key
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       fmt.Sprintf("%v", r.Id),
			ParentId: &parent,
			Name:     r.Name,
			FullName: r.Name,
			Data:     r.ConvertApiScope().(*models.BitbucketServerRepo),
		})
	}

	if resBody.IsLastPage {
		return children, nil, nil
	}

	if resBody.NextPageStart == nil {
		if len(resBody.Values) >= page.Limit {
			page.Start += len(resBody.Values)
		} else {
			return children, nil, nil
		}
	} else {
		if *resBody.NextPageStart < page.Start {
			return children, nil, nil
		}
		page.Start = *resBody.NextPageStart
	}

	return children, &page, nil
}

func searchBitbucketServerRepos(apiClient plugin.ApiClient, params *dsmodels.DsRemoteApiScopeSearchParams) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo],
	errors.Error,
) {
	query := initialQuery(BitBucketServerRemotePagination{
		Start: 0,
		Limit: 1,
	})
	query.Set("name", params.Search)

	// list repos part
	res, err := apiClient.Get("rest/api/latest/repos", query, nil)
	if err != nil {
		return nil, err
	}

	var resBody models.ReposResponse
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}

	children := []dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{}
	for _, r := range resBody.Values {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.BitbucketServerRepo]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       fmt.Sprintf("%v", r.Id),
			ParentId: nil,
			Name:     r.Name,
			FullName: r.Name,
			Data:     r.ConvertApiScope().(*models.BitbucketServerRepo),
		})
	}

	return children, nil
}

func initialQuery(page BitBucketServerRemotePagination) url.Values {
	query := url.Values{}
	query.Set("start", fmt.Sprintf("%v", page.Start))
	query.Set("limit", fmt.Sprintf("%v", page.Limit))
	return query
}

type BitBucketServerRemotePagination struct {
	Start int `json:"start"`
	Limit int `json:"limit"`
}
