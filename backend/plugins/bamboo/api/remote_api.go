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
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

type BambooRemotePagination struct {
	MaxResult  int `json:"max-result" validate:"required"`
	StartIndex int `json:"start-index" validate:"required"`
}

func listBambooRemoteScopes(
	connection *models.BambooConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page BambooRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.BambooPlan],
	nextPage *BambooRemotePagination,
	err errors.Error,
) {
	if page.MaxResult == 0 {
		page.MaxResult = 100
	}

	query := url.Values{
		"showEmpty":   []string{"false"},
		"max-result":  []string{fmt.Sprintf("%v", page.MaxResult)},
		"start-index": []string{fmt.Sprintf("%v", page.StartIndex)},
	}
	res, err := apiClient.Get("plan.json", query, nil)

	if err != nil {
		return
	}
	var planRes struct {
		Expand string `json:"expand"`
		Link   struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"link"`
		Plans struct {
			Size       int                    `json:"size"`
			Expand     string                 `json:"expand"`
			StartIndex int                    `json:"start-index"`
			MaxResult  int                    `json:"max-result"`
			Plan       []models.ApiBambooPlan `json:"plan"`
		} `json:"plans"`
	}
	err = api.UnmarshalResponse(res, &planRes)
	if err != nil {
		return
	}
	for _, plan := range planRes.Plans.Plan {
		children = append(children, toPlanModel(&plan))
	}
	// there may be more repos
	if len(children) == page.MaxResult {
		nextPage = &BambooRemotePagination{
			MaxResult:  page.MaxResult,
			StartIndex: page.StartIndex + page.MaxResult,
		}
	}
	return
}

func searchBambooPlans(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.BambooPlan],
	err errors.Error,
) {
	res, err := apiClient.Get(
		"search/plans.json",
		url.Values{
			"searchTerm": []string{params.Search},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resBody := models.ApiBambooSearchPlanResponse{}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	for _, r := range resBody.SearchResults {
		children = append(children, toPlanModel(&models.ApiBambooPlan{
			Key:  r.SearchEntity.Key,
			Name: r.SearchEntity.Name(),
		}))
	}
	return
}

func toPlanModel(plan *models.ApiBambooPlan) dsmodels.DsRemoteApiScopeListEntry[models.BambooPlan] {
	return dsmodels.DsRemoteApiScopeListEntry[models.BambooPlan]{
		Type:     api.RAS_ENTRY_TYPE_SCOPE,
		Id:       plan.Key,
		Name:     plan.Name,
		FullName: plan.Name,
		Data: &models.BambooPlan{
			PlanKey:   plan.Key,
			Name:      plan.Name,
			ShortName: plan.ShortName,
			ShortKey:  plan.ShortKey,
			Type:      plan.Type,
			Enabled:   plan.Enabled,
			Href:      plan.Link.Href,
			Rel:       plan.Link.Rel,
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
// @Tags plugins/bamboo
// @Router /plugins/bamboo/connections/{connectionId}/remote-scopes [GET]
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
// @Tags plugins/bamboo
// @Router /plugins/bamboo/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/bamboo
// @Router /plugins/bamboo/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
