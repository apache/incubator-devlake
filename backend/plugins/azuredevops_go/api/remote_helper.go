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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"net/url"
	"strings"
)

const (
	itemsPerPage = 100
	idSeparator  = "/"
)

func listAzuredevopsRemoteScopes(
	connection *models.AzuredevopsConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page AzuredevopsRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	nextPage *AzuredevopsRemotePagination,
	err errors.Error,
) {
	if page.Top == 0 {
		page.Top = itemsPerPage
	}

	if groupId != "" {
		id := strings.Split(groupId, idSeparator)
		return listAzuredevopsRepos(apiClient, id[0], id[1])
	}
	return listAzuredevopsProjects(connection, apiClient, page)
}

func listAzuredevopsProjects(
	connection *models.AzuredevopsConnection,
	apiClient plugin.ApiClient,
	page AzuredevopsRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	nextPage *AzuredevopsRemotePagination,
	err errors.Error) {

	query := url.Values{}
	query.Set("$top", fmt.Sprint(page.Top))
	query.Set("$skip", fmt.Sprint(page.Skip))
	query.Set("api-version", "7.1")

	vsc := newVsClient(connection, "https://app.vssps.visualstudio.com")

	profile, err := vsc.UserProfile()
	if err != nil {
		return nil, nil, err
	}
	accounts, err := vsc.UserAccounts(profile.Id)
	if err != nil {
		return nil, nil, err
	}

	var data struct {
		Projects []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsProject] `json:"value"`
	}

	for _, v := range accounts {
		res, err := apiClient.Get(fmt.Sprintf("%s/_apis/projects", v.AccountName), query, nil)
		if err != nil {
			return nil, nil, err
		}
		err = api.UnmarshalResponse(res, &data)
		if err != nil {
			if err != nil {
				return nil, nil, err
			}
		}

		for _, vv := range data.Projects {
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo]{
				Id:   v.AccountName + idSeparator + vv.Name,
				Type: api.RAS_ENTRY_TYPE_GROUP,
				Name: vv.Name,
			})
		}
	}

	if len(data.Projects) >= itemsPerPage {
		nextPage = &AzuredevopsRemotePagination{
			Top:  itemsPerPage,
			Skip: page.Skip + itemsPerPage,
		}
	}
	return
}

func listAzuredevopsRepos(
	apiClient plugin.ApiClient,
	orgId, projectId string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	nextPage *AzuredevopsRemotePagination,
	err errors.Error) {

	query := url.Values{}
	query.Set("api-version", "7.1")

	var data struct {
		Repos []AzuredevopsApiRepo `json:"value"`
	}

	res, err := apiClient.Get(fmt.Sprintf("%s/%s/_apis/git/repositories", orgId, projectId), query, nil)
	if err != nil {
		return nil, nil, err
	}
	err = api.UnmarshalResponse(res, &data)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range data.Repos {
		pID := orgId + idSeparator + projectId
		repo := v.toRepoModel()
		repo.ProjectId = projectId
		repo.OrganizationId = orgId
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: &pID,
			Id:       v.Id,
			Name:     v.Name,
			FullName: v.Name,
			Data:     &repo,
		})
	}
	return
}
