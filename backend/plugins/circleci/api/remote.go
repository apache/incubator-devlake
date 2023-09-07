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
	gocontext "context"
	"net/http"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
)

// RemoteApiScope
var _ plugin.ApiScope = (*RemoteProject)(nil)

type RemoteProject struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	OrgId string `json:"org_id"`
	Slug  string `json:"slug"`
}

// ConvertApiScope implements plugin.ApiScope.
func (c RemoteProject) ConvertApiScope() plugin.ToolLayerScope {
	return &models.CircleciProject{
		Id:             c.Id,
		Name:           c.Name,
		OrganizationId: c.OrgId,
		Slug:           c.Slug,
	}
}

func getRemoteProjects(
	basicRes context.BasicRes,
	gid string,
	queryData *api.RemoteQueryData,
	connection models.CircleciConnection,
) ([]RemoteProject, errors.Error) {
	// create api client
	apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("private/me", nil, http.Header{
		"Accept": []string{"application/json"},
	})
	if err != nil {
		return nil, err
	}

	var projects struct {
		FollowedProjects []RemoteProject `json:"followed_projects"`
	}
	err = api.UnmarshalResponse(res, &projects)
	if err != nil {
		return nil, err
	}
	// hack the queryData to stop pagination
	queryData.PerPage = len(projects.FollowedProjects) + 1
	return projects.FollowedProjects, err
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/circleci
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/circleci/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input, nil, getRemoteProjects)
}
