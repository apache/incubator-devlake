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
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

type ZentaoRemotePagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type ZentaoRemoteProjects struct {
	ZentaoRemotePagination
	Total  int                    `json:"total"`
	Values []models.ZentaoProject `json:"projects"`
}

func listZentaoRemoteScopes(
	connection *models.ZentaoConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page ZentaoRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.ZentaoProject],
	nextPage *ZentaoRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	if page.Limit == 0 {
		page.Limit = 20
	}
	// list projects part
	res, err := apiClient.Get("/projects", url.Values{
		"page":  {fmt.Sprintf("%d", page.Page)},
		"limit": {fmt.Sprintf("%d", page.Limit)},
	}, nil)
	if err != nil {
		return
	}
	// parse response body
	resBody := &ZentaoRemoteProjects{}
	err = api.UnmarshalResponse(res, resBody)
	if err != nil {
		return
	}
	// convert to dsmodels.DsRemoteApiScopeListEntry
	for _, p := range resBody.Values {
		tmpProject := p
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.ZentaoProject]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       fmt.Sprintf("%v", tmpProject.Id),
			Name:     tmpProject.Name,
			FullName: tmpProject.Name,
			Data:     &tmpProject,
		})
	}
	// next page
	if (resBody.Page-1)*resBody.Limit+len(resBody.Values) < resBody.Total {
		nextPage = &ZentaoRemotePagination{
			Page:  page.Page + 1,
			Limit: page.Limit,
		}
	}
	return
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/zentao
// @Router /plugins/zentao/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
