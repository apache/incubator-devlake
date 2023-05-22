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
	"strings"

	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BitbucketConnection) ([]models.GroupResponse, errors.Error) {
			if gid != "" {
				return nil, nil
			}
			query := initialQuery(queryData)

			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			query.Set("sort", "workspace.slug")
			query.Set("fields", "values.workspace.slug,values.workspace.name,pagelen,page,size")
			res, err = apiClient.Get("/user/permissions/workspaces", query, nil)
			if err != nil {
				return nil, err
			}

			resBody := &models.WorkspaceResponse{}
			err = api.UnmarshalResponse(res, resBody)
			if err != nil {
				return nil, err
			}

			return resBody.Values, err
		},
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BitbucketConnection) ([]models.BitbucketApiRepo, errors.Error) {
			if gid == "" {
				return nil, nil
			}
			query := initialQuery(queryData)

			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.display_name,values.created_on,values.updated_on,values.links.clone,values.links.html,pagelen,page,size")
			// list projects part
			res, err = apiClient.Get(fmt.Sprintf("/repositories/%s", gid), query, nil)
			if err != nil {
				return nil, err
			}
			var resBody models.ReposResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody.Values, err
		},
	)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/bitbucket
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context2.BasicRes, queryData *api.RemoteQueryData, connection models.BitbucketConnection) ([]models.BitbucketApiRepo, errors.Error) {
			// create api client
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, err
			}
			query := initialQuery(queryData)
			s := queryData.Search[0]

			// request search
			query.Set("sort", "name")
			query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.username,values.created_on,values.updated_on,values.links.clone,values.links.self,pagelen,page,size")
			gid := ``
			if strings.Contains(s, `/`) {
				gid = strings.Split(s, `/`)[0]
				s = strings.Split(s, `/`)[0]
			}
			query.Set("q", fmt.Sprintf(`name~"%s"`, s))
			// list repos part
			res, err := apiClient.Get(fmt.Sprintf("/repositories/%s", gid), query, nil)
			if err != nil {
				return nil, err
			}
			resBody := &models.ReposResponse{}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody.Values, err
		},
	)
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", queryData.Page))
	query.Set("pagelen", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
