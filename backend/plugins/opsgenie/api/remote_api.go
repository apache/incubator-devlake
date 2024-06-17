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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models/raw"
)

type OpsgenieRemotePagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type ServiceResponse struct {
	TotalCount int           `json:"totalCount"`
	Data       []raw.Service `json:"data"`
}

func listOpsgenieRemoteScopes(
	connection *models.OpsgenieConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page OpsgenieRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.Service],
	nextPage *OpsgenieRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	if page.PerPage == 0 {
		page.PerPage = 100
	}

	query := url.Values{
		"page":     []string{fmt.Sprintf("%v", page.Page)},
		"per_page": []string{fmt.Sprintf("%v", page.PerPage)},
	}

	res, err := apiClient.Get("v1/services", query, nil)
	if err != nil {
		return nil, nil, err
	}
	response := &ServiceResponse{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return nil, nil, err
	}

	// append service to output
	for _, service := range response.Data {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.Service]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       service.Id,
			Name:     service.Name,
			FullName: service.Name,
			Data: &models.Service{
				Url:    service.Links.Web,
				Id:     service.Id,
				Name:   service.Name,
				TeamId: service.TeamId,
				Scope: common.Scope{
					NoPKModel:    common.NoPKModel{},
					ConnectionId: connection.ID,
				},
			},
		})
	}

	return
}

// RemoteScopes list all available scopes (services) for this connection
// @Summary list all available scopes (services) for this connection
// @Description list all available scopes (services) for this connection
// @Tags plugins/opsgenie
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.Service]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/opsgenie/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/opsgenie
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.Service]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/opsgenie/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// Not supported
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusMethodNotAllowed}, nil
}
