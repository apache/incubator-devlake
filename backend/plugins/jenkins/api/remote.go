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
	getter := func(basicRes context.BasicRes, groupId string, queryData *api.RemoteQueryData, connection models.JenkinsConnection) (*api.RemoteScopesOutput, errors.Error) {
		var resBody []api.RemoteScopesChild
		apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
		}
		getJobsPageCallBack := func(job *models.Job) errors.Error {
			if job.Jobs != nil {
				// this is a group
				job.Path = groupId
				groupChild := api.RemoteScopesChild{
					Type: api.TypeGroup,
					Id:   job.GroupId(),
					Name: job.GroupName(),
					// don't need to save group into data
					Data: nil,
				}
				groupChild.ParentId = &groupId
				if *groupChild.ParentId == "" {
					groupChild.ParentId = nil
				}
				resBody = append(resBody, groupChild)
			} else {
				// this is a scope
				scope := job.ConvertApiScope()
				scopeChild := api.RemoteScopesChild{
					Type:     api.TypeScope,
					Id:       scope.ScopeId(),
					Name:     scope.ScopeName(),
					FullName: scope.ScopeFullName(),
					Data:     &scope,
				}
				scopeChild.ParentId = &groupId
				if *scopeChild.ParentId == "" {
					scopeChild.ParentId = nil
				}
				resBody = append(resBody, scopeChild)
			}
			return nil
		}
		_, err = GetJobsPage(apiClient, groupId, queryData.Page-1, queryData.PerPage, getJobsPageCallBack)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "failed to GetJobsPage")
		}
		return &api.RemoteScopesOutput{
			Children: resBody,
		}, nil
	}
	remoteScopesOutput, err := remoteHelper.GetRemoteScopesOutput(input, getter)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: remoteScopesOutput, Status: http.StatusOK}, nil
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
