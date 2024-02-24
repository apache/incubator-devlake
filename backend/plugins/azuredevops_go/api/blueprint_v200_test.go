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
	"testing"

	coreModels "github.com/apache/incubator-devlake/core/models"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/tasks"
	"github.com/stretchr/testify/assert"
)

const (
	connectionID        uint64 = 1
	azuredevopsRepoId          = "ad05901f-c9b0-4938-bc8a-a22eb2467ceb"
	expectDomainScopeId        = "azuredevops_go:AzuredevopsRepo:1:ad05901f-c9b0-4938-bc8a-a22eb2467ceb"
)

func mockAzuredevopsPlugin(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/azuredevops_go")
	mockMeta.On("Name").Return("dummy").Maybe()
	err := plugin.RegisterPlugin("azuredevops_go", mockMeta)
	assert.Equal(t, err, nil)
}

func TestMakeScopes(t *testing.T) {
	mockAzuredevopsPlugin(t)

	actualScopes, err := makeScopeV200(
		connectionID,
		[]*srvhelper.ScopeDetail[models.AzuredevopsRepo, models.AzuredevopsScopeConfig]{
			{
				Scope: models.AzuredevopsRepo{
					Scope: common.Scope{
						ConnectionId: connectionID,
					},
					Id: azuredevopsRepoId,
				},
				ScopeConfig: &models.AzuredevopsScopeConfig{
					ScopeConfig: common.ScopeConfig{
						Entities: []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CICD},
					},
				},
			},
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(actualScopes))
	assert.Equal(t, actualScopes[0].ScopeId(), expectDomainScopeId)
	assert.Equal(t, actualScopes[1].ScopeId(), expectDomainScopeId)
	assert.Equal(t, actualScopes[2].ScopeId(), expectDomainScopeId)
}

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockAzuredevopsPlugin(t)

	const (
		httpUrlToRepo          = "https://this_is_cloneUrl"
		azureDevOpsToken       = "personal-access-token"
		azureDevOpsProjectName = "azuredevops-test-project"
		azureDevOpsOrgName     = "azuredevops-test-org"
	)

	actualPlans, err := makePipelinePlanV200(
		[]plugin.SubTaskMeta{
			tasks.CollectApiPullRequestsMeta,
			tasks.ExtractApiPullRequestsMeta,
			tasks.CollectBuildsMeta,
			tasks.ExtractApiBuildsMeta,
		},
		&models.AzuredevopsConnection{
			BaseConnection: api.BaseConnection{
				Model: common.Model{
					ID: connectionID,
				},
			},
			AzuredevopsConn: models.AzuredevopsConn{
				AzuredevopsAccessToken: models.AzuredevopsAccessToken{
					Token: azureDevOpsToken,
				},
			},
		},
		[]*srvhelper.ScopeDetail[models.AzuredevopsRepo, models.AzuredevopsScopeConfig]{
			{
				Scope: models.AzuredevopsRepo{
					Id: fmt.Sprint(azuredevopsRepoId),
					AzureDevOpsPK: models.AzureDevOpsPK{
						ProjectId:      azureDevOpsProjectName,
						OrganizationId: azureDevOpsOrgName,
					},
					Name:      azureDevOpsProjectName,
					Url:       httpUrlToRepo,
					RemoteUrl: httpUrlToRepo,
				},
				ScopeConfig: &models.AzuredevopsScopeConfig{
					ScopeConfig: common.ScopeConfig{
						Entities: []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_CODE_REVIEW, plugin.DOMAIN_TYPE_CICD},
					},
					DeploymentPattern: "(?i)deploy",
					ProductionPattern: "(?i)prod",
					Refdiff: map[string]interface{}{
						"tagsPattern": "pattern",
						"tagsLimit":   10,
						"tagsOrder":   "reverse semver",
					},
				},
			},
		},
	)
	assert.Nil(t, err)

	var expectPlans = coreModels.PipelinePlan{
		{
			{
				Plugin: "azuredevops_go",
				Subtasks: []string{
					tasks.CollectApiPullRequestsMeta.Name,
					tasks.ExtractApiPullRequestsMeta.Name,
					tasks.CollectBuildsMeta.Name,
					tasks.ExtractApiBuildsMeta.Name,
				},
				Options: map[string]interface{}{
					"connectionId":   connectionID,
					"projectId":      azureDevOpsProjectName,
					"repositoryId":   fmt.Sprint(azuredevopsRepoId),
					"organizationId": azureDevOpsOrgName,
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": expectDomainScopeId,
					"name":   azureDevOpsProjectName,
					"url":    "https://git:personal-access-token@this_is_cloneUrl",
				},
			},
		},
		{
			{
				Plugin: "refdiff",
				Options: map[string]interface{}{
					"tagsLimit":   10,
					"tagsOrder":   "reverse semver",
					"tagsPattern": "pattern",
				},
			},
		},
	}

	assert.Equal(t, expectPlans, actualPlans)
}
