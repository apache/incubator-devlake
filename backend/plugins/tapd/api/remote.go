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
	"encoding/json"
	"fmt"
	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
	"io"
	"net/http"
	"net/url"
)

// PrepareFirstPageToken prepare first page token
// @Summary prepare first page token
// @Description prepare first page token
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param companyId query string false "company ID"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/remote-scopes-prepare-token [GET]
func PrepareFirstPageToken(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.PrepareFirstPageToken(input.Query[`companyId`][0])
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.TapdConnection) ([]api.BaseRemoteGroupResponse, errors.Error) {
			if gid == "" {
				// if gid is empty, it means we need to query company
				gid = "1"
			}
			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			query := url.Values{}
			query.Set("company_id", queryData.CustomInfo)
			res, err = apiClient.Get("/workspaces/projects", query, nil)
			if err != nil {
				return nil, err
			}

			var resBody models.WorkspacesResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			if resBody.Status != 1 {
				return nil, errors.BadInput.Wrap(err, "failed to get workspaces")
			}

			// check if workspace is a group
			isGroupMap := map[uint64]bool{}
			for _, workspace := range resBody.Data {
				isGroupMap[workspace.ApiTapdWorkspace.ParentId] = true
			}

			groups := []api.BaseRemoteGroupResponse{}
			for _, workspace := range resBody.Data {
				if fmt.Sprintf(`%d`, workspace.ApiTapdWorkspace.ParentId) == gid &&
					isGroupMap[workspace.ApiTapdWorkspace.Id] {
					groups = append(groups, api.BaseRemoteGroupResponse{
						Id:   fmt.Sprintf(`%d`, workspace.ApiTapdWorkspace.Id),
						Name: workspace.ApiTapdWorkspace.Name,
					})
				}
			}

			return groups, err
		},
		func(basicRes context2.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.TapdConnection) ([]models.TapdWorkspace, errors.Error) {
			if gid == "" {
				return nil, nil
			}

			apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			query := url.Values{}
			query.Set("company_id", queryData.CustomInfo)
			res, err = apiClient.Get("/workspaces/projects", query, nil)
			if err != nil {
				return nil, err
			}
			var resBody models.WorkspacesResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			workspaces := []models.TapdWorkspace{}
			for _, workspace := range resBody.Data {
				if fmt.Sprintf(`%d`, workspace.ApiTapdWorkspace.ParentId) == gid {
					// filter from all project to query what we need...
					workspaces = append(workspaces, models.TapdWorkspace(workspace.ApiTapdWorkspace))
				}

			}
			return workspaces, err
		},
	)
}

func GetApiWorkspace(op *tasks.TapdOptions, apiClient aha.ApiClientAbstract) (*models.TapdWorkspace, errors.Error) {
	query := url.Values{}
	query.Set("workspace_id", fmt.Sprintf("%v", op.WorkspaceId))
	res, err := apiClient.Get("workspaces/get_workspace_info", query, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting repo detail from %s", res.Request.URL.String()))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}

	var resBody models.WorkspaceResponse
	err = errors.Convert(json.Unmarshal(body, &resBody))
	if err != nil {
		return nil, err
	}
	workspace := models.TapdWorkspace(resBody.Data.ApiTapdWorkspace)
	workspace.ConnectionId = op.ConnectionId
	return &workspace, nil
}
