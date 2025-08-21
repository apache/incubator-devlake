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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/api/azuredevops"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"sync"
)

const (
	idSeparator    = "/"
	maxConcurrency = 10
)

type AzuredevopsRemotePagination struct {
	Skip int
	Top  int
}

// https://learn.microsoft.com/en-us/azure/devops/pipelines/repos/?view=azure-devops
// https://learn.microsoft.com/en-us/azure/devops/pipelines/repos/multi-repo-checkout?view=azure-devops
var supportedSourceRepositories = []string{"github", "githubenterprise", "bitbucket", "git"}

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

	org := connection.Organization
	vsc := azuredevops.NewClient(connection, apiClient, "https://app.vssps.visualstudio.com")

	if groupId == "" {
		return listAzuredevopsProjects(vsc, page, org)
	}

	id := strings.Split(groupId, idSeparator)

	if remote, err := listRemoteRepos(vsc, id[0], id[1]); err == nil {
		children = append(children, remote...)
	}

	if remote, err := listAzuredevopsRepos(vsc, id[0], id[1]); err == nil {
		children = append(children, remote...)
	}
	return children, nextPage, nil
}

func listAzuredevopsProjects(vsc azuredevops.Client, _ AzuredevopsRemotePagination, org string) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	nextPage *AzuredevopsRemotePagination,
	err errors.Error) {

	var accounts azuredevops.AccountResponse
	if org == "" {
		profile, err := vsc.GetUserProfile()
		if err != nil {
			return nil, nil, err
		}
		accounts, err = vsc.GetUserAccounts(profile.Id)
		if err != nil {
			return nil, nil, err
		}
	} else {
		accounts = append(accounts, azuredevops.Account{AccountName: org})
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(maxConcurrency)

	var mu sync.Mutex

	for _, v := range accounts {
		accountName := v.AccountName
		g.Go(func() error {
			args := azuredevops.GetProjectsArgs{
				OrgId: accountName,
			}
			projects, err := vsc.GetProjects(args)
			if err != nil {
				return err
			}

			var tmp []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo]
			for _, vv := range projects {
				tmp = append(tmp, dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo]{
					Id:   accountName + idSeparator + vv.Name,
					Type: api.RAS_ENTRY_TYPE_GROUP,
					Name: vv.Name,
				})
			}
			mu.Lock()
			children = append(children, tmp...)
			mu.Unlock()
			return nil
		})
	}

	err = errors.Convert(g.Wait())
	return
}

func listAzuredevopsRepos(
	vsc azuredevops.Client,
	orgId, projectId string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	err errors.Error) {

	args := azuredevops.GetRepositoriesArgs{
		OrgId:     orgId,
		ProjectId: projectId,
	}

	repos, err := vsc.GetRepositories(args)
	if err != nil {
		return nil, err
	}

	for _, v := range repos {
		if v.IsDisabled {
			continue
		}

		pID := orgId + idSeparator + projectId
		repo := models.AzuredevopsRepo{
			Id:        v.Id,
			Type:      models.RepositoryTypeADO,
			Name:      v.Name,
			Url:       v.Url,
			RemoteUrl: v.RemoteUrl,
			IsFork:    false,
		}
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

func listRemoteRepos(
	vsc azuredevops.Client,
	orgId, projectId string,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo],
	err errors.Error) {

	args := azuredevops.GetServiceEndpointsArgs{
		OrgId:     orgId,
		ProjectId: projectId,
	}

	endpoints, err := vsc.GetServiceEndpoints(args)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var remoteRepos []azuredevops.RemoteRepository

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(maxConcurrency)

	for _, v := range endpoints {
		if !slices.Contains(supportedSourceRepositories, v.Type) {
			continue
		}

		remoteRepoArgs := azuredevops.GetRemoteRepositoriesArgs{
			ProjectId:       projectId,
			OrgId:           orgId,
			Provider:        v.Type,
			ServiceEndpoint: v.Id,
		}

		g.Go(func() error {
			repos, err := vsc.GetRemoteRepositories(remoteRepoArgs)
			mu.Lock()
			remoteRepos = append(remoteRepos, repos...)
			mu.Unlock()
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Internal.Wrap(err, "failed to call 'GetRemoteRepositories', falling back to empty list")
	}

	for _, v := range remoteRepos {
		pID := orgId + idSeparator + projectId
		isFork, _ := strconv.ParseBool(v.Properties.IsFork)
		isPrivate, _ := strconv.ParseBool(v.Properties.IsPrivate)

		// IDs must not contain URL reserved characters (e.g., "/"), as this breaks the routing in the scope API.
		// Accessing /plugins/azuredevops_go/connections/<id>/apache/incubator-devlake results in a 404 error, where
		// "apache/incubator-devlake" is the repository ID returned by ADOs sourceProviders API.
		// Therefore, we are creating our own ID, by combining the Service Connection and the External ID
		remoteId := fmt.Sprintf("%s-%s", v.Properties.ConnectedServiceId, v.Properties.ExternalId)

		repo := models.AzuredevopsRepo{
			Id:         remoteId,
			Type:       v.SourceProviderName,
			Name:       v.SourceProviderName + idSeparator + v.FullName,
			Url:        v.Properties.ManageUrl,
			RemoteUrl:  v.Properties.CloneUrl,
			ExternalId: v.Id,
			IsFork:     isFork,
			IsPrivate:  isPrivate,
		}

		repo.ProjectId = projectId
		repo.OrganizationId = orgId
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.AzuredevopsRepo]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: &pID,
			Id:       v.Id,
			Name:     v.Name,
			FullName: v.FullName,
			Data:     &repo,
		})
	}
	return
}
