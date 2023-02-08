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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	connection := &models.GithubConnection{
		BaseConnection: helper.BaseConnection{
			Name: "github-test",
			Model: common.Model{
				ID: 1,
			},
		},
		GithubConn: models.GithubConn{
			RestConnection: helper.RestConnection{
				Endpoint:         "https://api.github.com/",
				Proxy:            "",
				RateLimitPerHour: 0,
			},
			GithubAccessToken: models.GithubAccessToken{
				AccessToken: helper.AccessToken{
					Token: "123",
				},
			},
		},
	}
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")
	err := plugin.RegisterPlugin("github", mockMeta)
	assert.Nil(t, err)
	// Refresh Global Variables and set the sql mock
	basicRes = NewMockBasicRes()
	bs := &plugin.BlueprintScopeV200{
		Entities: []string{"CODE", "TICKET"},
		Id:       "1",
	}
	bpScopes := make([]*plugin.BlueprintScopeV200, 0)
	bpScopes = append(bpScopes, bs)
	syncPolicy := &plugin.BlueprintSyncPolicy{}

	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(nil, plan, bpScopes, connection, syncPolicy)
	assert.Nil(t, err)
	basicRes = NewMockBasicRes()
	scopes, err := makeScopesV200(bpScopes, connection)
	assert.Nil(t, err)

	expectPlan := plugin.PipelinePlan{
		plugin.PipelineStage{
			{
				Plugin:   "github",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"githubId":     12345,
					"name":         "test/testRepo",
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": "github:GithubRepo:1:12345",
					"url":    "https://git:123@this_is_cloneUrl",
				},
			},
		},
		plugin.PipelineStage{
			{
				Plugin: "refdiff",
				Options: map[string]interface{}{
					"repoId":      "github:GithubRepo:1:12345",
					"tagsLimit":   10,
					"tagsOrder":   "reverse semver",
					"tagsPattern": "pattern",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)
	expectScopes := make([]plugin.Scope, 0)
	scopeRepo := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: "github:GithubRepo:1:12345",
		},
		Name: "test/testRepo",
	}

	scopeTicket := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "github:GithubRepo:1:12345",
		},
		Name:        "test/testRepo",
		Description: "",
		Url:         "",
		CreatedDate: nil,
		Type:        "",
	}

	expectScopes = append(expectScopes, scopeRepo, scopeTicket)
	assert.Equal(t, expectScopes, scopes)
}

// NewMockBasicRes FIXME ...
func NewMockBasicRes() *mockcontext.BasicRes {
	testGithubRepo := &models.GithubRepo{
		ConnectionId:         1,
		GithubId:             12345,
		Name:                 "test/testRepo",
		CloneUrl:             "https://this_is_cloneUrl",
		TransformationRuleId: 1,
	}

	testTransformationRule := &models.GithubTransformationRule{
		Model: common.Model{
			ID: 1,
		},
		Name:   "github transformation rule",
		PrType: "hey,man,wasup",
		Refdiff: map[string]interface{}{
			"tagsPattern": "pattern",
			"tagsLimit":   10,
			"tagsOrder":   "reverse semver",
		},
	}
	mockRes := new(mockcontext.BasicRes)
	mockDal := new(mockdal.Dal)

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.GithubRepo)
		*dst = *testGithubRepo
	}).Return(nil).Once()

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.GithubTransformationRule)
		*dst = *testTransformationRule
	}).Return(nil).Once()

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")

	return mockRes
}
