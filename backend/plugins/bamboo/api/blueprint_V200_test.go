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

	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	const testConnectionID uint64 = 1
	const testTransformationRuleId uint64 = 2
	const testKey string = "TEST"
	const testBambooEndPoint string = "http://mail.nddtf.com:8085/rest/api/latest/"
	const testLink string = "http://mail.nddtf.com:8085/rest/api/latest/project/TEST"
	const testUser string = "username"
	const testPass string = "password"
	const testName string = "bamboo-test"
	const testTransformationRuleName string = "bamboo transformation rule"
	const testProxy string = ""

	syncPolicy := &plugin.BlueprintSyncPolicy{}
	bpScopes := []*plugin.BlueprintScopeV200{
		{
			Entities: []string{plugin.DOMAIN_TYPE_CICD},
			Id:       testKey,
			Name:     testName,
		},
	}

	var testBambooProject = &models.BambooProject{
		ConnectionId: testConnectionID,
		ProjectKey:   testKey,
		Name:         testName,
		Href:         testLink,

		TransformationRuleId: testTransformationRuleId,
	}

	var testTransformationRule = &models.BambooTransformationRule{
		Model: common.Model{
			ID: testTransformationRuleId,
		},
		Name: testTransformationRuleName,
	}

	var testBambooConnection = &models.BambooConnection{
		BaseConnection: helper.BaseConnection{
			Name: testName,
			Model: common.Model{
				ID: testConnectionID,
			},
		},
		BambooConn: models.BambooConn{
			RestConnection: helper.RestConnection{
				Endpoint:         testBambooEndPoint,
				Proxy:            testProxy,
				RateLimitPerHour: 0,
			},
			BasicAuth: helper.BasicAuth{
				Username: testUser,
				Password: testPass,
			},
		},
	}

	var expectRepoId = "bamboo:BambooProject:1:TEST"

	var testSubTaskMeta = []plugin.SubTaskMeta{
		tasks.ConvertProjectsMeta,
	}

	var expectPlans = plugin.PipelinePlan{
		{
			{
				Plugin: "bamboo",
				Subtasks: []string{
					tasks.ConvertProjectsMeta.Name,
				},
				Options: map[string]interface{}{
					"connectionId":         uint64(1),
					"projectKey":           testKey,
					"transformationRuleId": testTransformationRuleId,
				},
			},
		},
	}

	expectCicdScope := devops.NewCicdScope(expectRepoId, testName)
	expectCicdScope.Description = ""
	expectCicdScope.Url = ""

	var err errors.Error

	// register bamboo plugin for NewDomainIdGenerator
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/bamboo")
	err = plugin.RegisterPlugin("bamboo", mockMeta)
	assert.Equal(t, err, nil)

	// Refresh Global Variables and set the sql mock
	basicRes = unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.BambooConnection)
			*dst = *testBambooConnection
		}).Return(nil).Once()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.BambooProject)
			*dst = *testBambooProject
		}).Return(nil).Twice()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.BambooTransformationRule)
			*dst = *testTransformationRule
		}).Return(nil).Once()
	})
	connectionHelper = helper.NewConnectionHelper(
		basicRes,
		validator.New(),
	)

	plans, scopes, err := MakePipelinePlanV200(testSubTaskMeta, testConnectionID, bpScopes, syncPolicy)
	assert.Equal(t, err, nil)

	assert.Equal(t, expectPlans, plans)

	// ignore CreatedDate UpdatedDate  CreatedAt UpdatedAt checking
	expectCicdScope.CreatedDate = scopes[0].(*devops.CicdScope).CreatedDate
	expectCicdScope.UpdatedDate = scopes[0].(*devops.CicdScope).UpdatedDate
	expectCicdScope.CreatedAt = scopes[0].(*devops.CicdScope).CreatedAt
	expectCicdScope.UpdatedAt = scopes[0].(*devops.CicdScope).UpdatedAt

	var expectScopes = []plugin.Scope{
		expectCicdScope,
	}

	assert.Equal(t, expectScopes, scopes)
}
