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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/jenkins")
	err := plugin.RegisterPlugin("jenkins", mockMeta)
	assert.Nil(t, err)
	bs := &plugin.BlueprintScopeV200{
		Id:   "a/b/ccc",
		Name: "",
	}
	syncPolicy := &plugin.BlueprintSyncPolicy{}
	bpScopes := make([]*plugin.BlueprintScopeV200, 0)
	bpScopes = append(bpScopes, bs)

	basicRes = NewMockBasicRes()
	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(basicRes, nil, plan, bpScopes, 1, syncPolicy)
	assert.Nil(t, err)
	basicRes = NewMockBasicRes()
	scopes, err := makeScopesV200(basicRes, bpScopes, 1)
	assert.Nil(t, err)

	expectPlan := plugin.PipelinePlan{
		plugin.PipelineStage{
			{
				Plugin:   "jenkins",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"scopeId":      "a/b/ccc",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)

	expectScopes := make([]plugin.Scope, 0)
	scopeCicd := &devops.CicdScope{
		DomainEntity: domainlayer.DomainEntity{
			Id: "jenkins:JenkinsJob:1:a/b/ccc",
		},
		Name: "a/b/ccc",
	}

	expectScopes = append(expectScopes, scopeCicd)
	assert.Equal(t, expectScopes, scopes)
}

// NewMockBasicRes FIXME ...
func NewMockBasicRes() *mockcontext.BasicRes {
	jenkinsJob := &models.JenkinsJob{
		ConnectionId: 1,
		FullName:     "a/b/ccc",
	}

	scopeConfig := &models.JenkinsScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Entities: []string{"CICD"},
		},
	}

	mockRes := new(mockcontext.BasicRes)
	mockDal := new(mockdal.Dal)

	mockDal.On("First", mock.AnythingOfType("*models.JenkinsScopeConfig"), mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.JenkinsScopeConfig)
		*dst = *scopeConfig
	}).Return(nil)
	mockDal.On("First", mock.AnythingOfType("*models.JenkinsJob"), mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.JenkinsJob)
		*dst = *jenkinsJob
	}).Return(nil)

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")

	return mockRes
}
