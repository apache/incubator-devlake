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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/bitbucket_server")
	mockMeta.On("Name").Return("bitbucket_server").Maybe()
	err := plugin.RegisterPlugin("bitbucket_server", mockMeta)
	assert.Nil(t, err)
	// Refresh Global Variables and set the sql mock
	mockBasicRes(t)

	connection := &models.BitbucketServerConnection{
		BaseConnection: helper.BaseConnection{
			Name: "TP/repos/first-repo",
			Model: common.Model{
				ID: 3,
			},
		},
		BitbucketServerConn: models.BitbucketServerConn{
			RestConnection: helper.RestConnection{
				Endpoint:         "https://api.bitbucket.org/2.0/",
				Proxy:            "",
				RateLimitPerHour: 0,
			},
			BasicAuth: helper.BasicAuth{
				Username: "Username",
				Password: "Password",
			},
		},
	}

	scopeDetails := make([]*srvhelper.ScopeDetail[models.BitbucketServerRepo, models.BitbucketServerScopeConfig], 0)
	scopeDetails = append(scopeDetails, &srvhelper.ScopeDetail[models.BitbucketServerRepo, models.BitbucketServerScopeConfig]{
		Scope: models.BitbucketServerRepo{
			BitbucketId: "3",
		},
	})

	plan, err := makeDataSourcePipelinePlanV200(nil, scopeDetails, connection)
	assert.Nil(t, err)
	scopes, err := makeScopesV200(scopeDetails, connection)
	assert.Nil(t, err)

	expectPlan := coreModels.PipelinePlan{
		coreModels.PipelineStage{
			{
				Plugin:   "bitbucket_server",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"fullName":     "TP/repos/first-repo",
					"connectionId": uint64(3),
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": "bitbucket_server:BitbucketServerRepo:3:TP/repos/first-repo",
					"name":   "TP/repos/first-repo",
					"url":    "https://Username:Password@this_is_cloneUrl",
				},
			},
		},
		coreModels.PipelineStage{
			{
				Plugin: "refdiff",
				Options: map[string]interface{}{
					"repoId":      "bitbucket_server:BitbucketServerRepo:3:TP/repos/first-repo",
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
			Id: "bitbucket_server:BitbucketServerRepo:3:TP/repos/first-repo",
		},
		Name: "TP/repos/first-repo",
	}

	scopeTicket := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "bitbucket_server:BitbucketServerRepo:3:TP/repos/first-repo",
		},
		Name:        "TP/repos/first-repo",
		Description: "",
		Url:         "",
		CreatedDate: nil,
		Type:        "",
	}

	expectScopes = append(expectScopes, scopeRepo, scopeTicket)
	assert.Equal(t, expectScopes, scopes)
}

func mockBasicRes(t *testing.T) {
	testBitbucketRepo := &models.BitbucketServerRepo{
		Scope: common.Scope{
			ConnectionId:  3,
			ScopeConfigId: 1,
		},
		BitbucketId: "TP/repos/first-repo",
		Name:        "test/testRepo",
		CloneUrl:    "https://this_is_cloneUrl",
	}

	testScopeConfig := &models.BitbucketServerScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Entities: []string{"CODE", "TICKET"},
			Model: common.Model{
				ID: 1,
			},
		},
		Name: "Bitbucket scope config",
		Refdiff: map[string]interface{}{
			"tagsPattern": "pattern",
			"tagsLimit":   10,
			"tagsOrder":   "reverse semver",
		},
	}
	// Refresh Global Variables and set the sql mock
	mockRes := unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.AnythingOfType("*models.BitbucketServerRepo"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.BitbucketServerRepo)
			*dst = *testBitbucketRepo
		}).Return(nil)

		mockDal.On("First", mock.AnythingOfType("*models.BitbucketServerScopeConfig"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.BitbucketServerScopeConfig)
			*dst = *testScopeConfig
		}).Return(nil)
	})
	p := mockplugin.NewPluginMeta(t)
	p.On("Name").Return("dummy").Maybe()
	Init(mockRes, p)
}
