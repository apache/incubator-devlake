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
	"testing"

	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/jira")
	mockMeta.On("Name").Return("jira").Maybe()
	err := plugin.RegisterPlugin("jira", mockMeta)
	assert.Nil(t, err)
	bs := &coreModels.BlueprintScope{
		ScopeId: "10",
	}
	syncPolicy := &coreModels.BlueprintSyncPolicy{}
	bpScopes := make([]*coreModels.BlueprintScope, 0)
	bpScopes = append(bpScopes, bs)
	plan := make(coreModels.PipelinePlan, len(bpScopes))
	mockBasicRes(t)

	plan, err = makeDataSourcePipelinePlanV200(nil, plan, bpScopes, uint64(1), syncPolicy)
	assert.Nil(t, err)
	scopes, err := makeScopesV200(bpScopes, uint64(1))
	assert.Nil(t, err)

	expectPlan := coreModels.PipelinePlan{
		coreModels.PipelineStage{
			{
				Plugin:   "jira",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"scopeId":      "10",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)

	expectScopes := make([]plugin.Scope, 0)
	jiraBoard := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "jira:JiraBoard:1:10",
		},
		Name: "a",
	}

	expectScopes = append(expectScopes, jiraBoard)
	assert.Equal(t, expectScopes, scopes)
}

func mockBasicRes(t *testing.T) {
	jiraBoard := &models.JiraBoard{
		ConnectionId: 1,
		BoardId:      10,
		Name:         "a",
	}
	scopeConfig := &models.JiraScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Entities: []string{plugin.DOMAIN_TYPE_TICKET},
		},
	}

	// Refresh Global Variables and set the sql mock
	mockRes := unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.AnythingOfType("*models.JiraScopeConfig"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.JiraScopeConfig)
			*dst = *scopeConfig
		}).Return(nil)
		mockDal.On("First", mock.AnythingOfType("*models.JiraBoard"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.JiraBoard)
			*dst = *jiraBoard
		}).Return(nil)
	})
	p := mockplugin.NewPluginMeta(t)
	p.On("Name").Return("dummy").Maybe()
	Init(mockRes, p)
}
