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
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	connection := &models.ZentaoConnection{
		BaseConnection: helper.BaseConnection{
			Name: "zentao-test",
			Model: common.Model{
				ID: 1,
			},
		},
		ZentaoConn: models.ZentaoConn{
			RestConnection: helper.RestConnection{
				Endpoint:         "https://zentao.example.org/api.php/v1/",
				Proxy:            "",
				RateLimitPerHour: 0,
			},
			BasicAuth: helper.BasicAuth{
				Username: "Username",
				Password: "Password",
			},
		},
	}
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/zentao")
	err := plugin.RegisterPlugin("zentao", mockMeta)
	assert.Nil(t, err)
	// Refresh Global Variables and set the sql mock
	basicRes = NewMockBasicRes()
	bs := &plugin.BlueprintScopeV200{
		Entities: []string{"TICKET"},
		Id:       "project/1",
	}
	bs2 := &plugin.BlueprintScopeV200{
		Entities: []string{"TICKET"},
		Id:       "product/1",
	}
	bpScopes := make([]*plugin.BlueprintScopeV200, 0)
	bpScopes = append(bpScopes, bs, bs2)
	syncPolicy := &plugin.BlueprintSyncPolicy{}

	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, scopes, err := makePipelinePlanV200(nil, plan, bpScopes, connection, syncPolicy)
	assert.Nil(t, err)
	basicRes = NewMockBasicRes()

	expectPlan := plugin.PipelinePlan{
		plugin.PipelineStage{
			{
				Plugin:   "zentao",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"ConnectionId": uint64(1),
					"productId":    int64(0),
					"projectId":    int64(1),
				},
			},
		},
		plugin.PipelineStage{
			{
				Plugin:   "zentao",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"ConnectionId": uint64(1),
					"productId":    int64(1),
					"projectId":    int64(0),
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)
	expectScopes := make([]plugin.Scope, 0)
	scopeTicket1 := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "zentao:ZentaoProject:1:1",
		},
		Name:        "test/testRepo",
		Description: "",
		Url:         "",
		CreatedDate: nil,
		Type:        `project`,
	}
	scopeTicket2 := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "zentao:ZentaoProduct:1:1",
		},
		Name:        "test/testRepo",
		Description: "",
		Url:         "",
		CreatedDate: nil,
		Type:        `product/normal`,
	}

	expectScopes = append(expectScopes, scopeTicket1, scopeTicket2)
	assert.Equal(t, expectScopes, scopes)
}

// NewMockBasicRes FIXME ...
func NewMockBasicRes() *mockcontext.BasicRes {
	testZentaoProduct := &models.ZentaoProduct{
		ConnectionId: 1,
		Id:           1,
		Name:         "test/testRepo",
		Type:         `product/normal`,
		//TransformationRuleId: 1,
	}
	testZentaoProject := &models.ZentaoProject{
		ConnectionId: 1,
		Id:           1,
		Name:         "test/testRepo",
		Type:         `project`,
		//TransformationRuleId: 1,
	}

	//testTransformationRule := &models.ZentaoTransformation{
	//	Model: common.Model{
	//		ID: 1,
	//	},
	//	Name:            "Zentao transformation rule",
	//}
	mockRes := new(mockcontext.BasicRes)
	mockDal := new(mockdal.Dal)

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.ZentaoProject)
		*dst = *testZentaoProject
	}).Return(nil).Once()

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.ZentaoProduct)
		*dst = *testZentaoProduct
	}).Return(nil).Once()

	//mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
	//	dst := args.Get(0).(*models.ZentaoTransformation)
	//	*dst = *testTransformationRule
	//}).Return(nil).Once()

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")

	return mockRes
}
