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
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

type ArgocdRemotePagination struct {
	Page int `json:"page"`
}

type ArgocdRemoteApplication struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Project   string `json:"project"`
}

func listArgocdRemoteScopes(
	connection *models.ArgocdConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page ArgocdRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.ArgocdApplication],
	nextPage *ArgocdRemotePagination,
	err errors.Error,
) {

	if groupId == "" {
		query := url.Values{}
		res, err := apiClient.Get("applications", query, nil)
		if err != nil {
			return nil, nextPage, err
		}

		var argocdApps struct {
			Items []struct {
				Metadata struct {
					Name      string `json:"name"`
					Namespace string `json:"namespace"`
				} `json:"metadata"`
				Spec struct {
					Project string `json:"project"`
				} `json:"spec"`
			} `json:"items"`
		}

		if err = api.UnmarshalResponse(res, &argocdApps); err != nil {
			return nil, nextPage, err
		}

		projectMap := make(map[string]bool)
		for _, item := range argocdApps.Items {
			projectName := item.Spec.Project
			if projectName == "" {
				projectName = "default"
			}
			projectMap[projectName] = true
		}

		// Return projects as groups
		for projectName := range projectMap {
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.ArgocdApplication]{
				Type:     api.RAS_ENTRY_TYPE_GROUP,
				Id:       projectName,
				Name:     projectName,
				FullName: projectName,
			})
		}

		return children, nextPage, nil
	}

	query := url.Values{}
	res, err := apiClient.Get("applications", query, nil)
	if err != nil {
		return nil, nextPage, err
	}

	var argocdApps struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				Project string `json:"project"`
			} `json:"spec"`
		} `json:"items"`
	}

	if err = api.UnmarshalResponse(res, &argocdApps); err != nil {
		return nil, nextPage, err
	}

	for _, item := range argocdApps.Items {
		projectName := item.Spec.Project
		if projectName == "" {
			projectName = "default"
		}

		if projectName == groupId {
			app := &models.ArgocdApplication{
				Name:      item.Metadata.Name,
				Namespace: item.Metadata.Namespace,
				Project:   item.Spec.Project,
			}
			app.ConnectionId = connection.ID

			parentId := groupId
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.ArgocdApplication]{
				Type:     api.RAS_ENTRY_TYPE_SCOPE,
				ParentId: &parentId,
				Id:       item.Metadata.Name,
				Name:     item.Metadata.Name,
				FullName: item.Metadata.Name,
				Data:     app,
			})
		}
	}

	return children, nextPage, nil
}

// RemoteScopes list all available applications for the given connection
// @Summary list all available applications for the given connection
// @Description list all available applications for the given connection
// @Tags plugins/argocd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.ArgocdApplication]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the ArgoCD server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/argocd
// @Router /plugins/argocd/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
