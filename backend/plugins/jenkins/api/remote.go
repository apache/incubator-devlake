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
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/jenkins
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetMixedGroupAndScopesFromRemote(
		input,
		func(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.JenkinsConnection) ([]models.Job, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var resBody []models.Job
			_, err = GetJobsPage(apiClient, gid, queryData.Page-1, queryData.PerPage, func(job *models.Job) errors.Error {
				job.Path = gid
				resBody = append(resBody, *job)
				return nil
			})
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to GetJobsPage")
			}
			return resBody, err
		},
	)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/jenkins
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context.BasicRes, queryData *api.RemoteQueryData, connection models.JenkinsConnection) ([]models.Job, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}

			var resBody []models.Job
			breakError := errors.Default.New("we need break from get all jobs for page full")
			count := 0
			pageoOffset := (queryData.Page - 1) * queryData.PerPage
			err = GetAllJobs(apiClient, "", "", queryData.PerPage, func(job *models.Job, isPath bool) errors.Error {
				if job.Jobs == nil {
					if strings.Contains(job.FullName, queryData.Search[0]) {
						if count >= pageoOffset {
							resBody = append(resBody, *job)
						}
						count++
					}
					if len(resBody) > queryData.PerPage {
						return breakError
					}
				}
				return nil
			})
			if err != nil && err != breakError {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}

			return resBody, nil
		})
}
