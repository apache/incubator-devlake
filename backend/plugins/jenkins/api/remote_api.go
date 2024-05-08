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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

type JenkinsRemotePagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func listJenkinsRemoteScopes(
	connection *models.JenkinsConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page JenkinsRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.JenkinsJob],
	nextPage *JenkinsRemotePagination,
	err errors.Error,
) {
	if page.Page == 0 {
		page.Page = 1
	}
	if page.PerPage == 0 {
		page.PerPage = 100
	}
	var parentId *string
	if groupId != "" {
		parentId = &groupId
	}
	getJobsPageCallBack := func(job *models.Job) errors.Error {
		switch job.Class {
		case "org.jenkinsci.plugins.workflow.job.WorkflowJob":
			fallthrough
		case "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject":
			fallthrough
		case "hudson.model.FreeStyleProject":
			// this is a scope
			jenkinsJob := job.ToJenkinsJob()
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.JenkinsJob]{
				Type:     api.RAS_ENTRY_TYPE_SCOPE,
				Id:       jenkinsJob.ScopeId(),
				Name:     jenkinsJob.ScopeName(),
				FullName: jenkinsJob.ScopeFullName(),
				Data:     jenkinsJob,
				ParentId: parentId,
			})
		default:
			// this is a group
			job.Path = groupId
			children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.JenkinsJob]{
				Type:     api.RAS_ENTRY_TYPE_GROUP,
				Id:       fmt.Sprintf("%s/job/%s", job.Path, job.Name),
				Name:     job.Name,
				ParentId: parentId,
			})
		}

		return nil
	}
	_, err = GetJobsPage(apiClient, groupId, page.Page-1, page.PerPage, getJobsPageCallBack)
	if err != nil {
		return
	}
	if len(children) == page.PerPage {
		nextPage = &JenkinsRemotePagination{
			Page:    page.Page + 1,
			PerPage: page.PerPage,
		}
	}
	return
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
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.JenkinsJob]
// @Tags plugins/jenkins
// @Router /plugins/jenkins/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Router /plugins/jenkins/connections/{connectionId}/proxy/{path} [GET]
// @Tags plugins/jenkins
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
