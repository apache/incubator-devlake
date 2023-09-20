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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/jenkins")
	mockMeta.On("Name").Return("jenkins").Maybe()
	err := plugin.RegisterPlugin("jenkins", mockMeta)
	assert.Nil(t, err)
	bs := &coreModels.BlueprintScope{
		ScopeId: "a/b/ccc",
	}
	
	bpScopes := make([]*coreModels.BlueprintScope, 0)
	bpScopes = append(bpScopes, bs)

	mockBasicRes(t)

	plan := make(coreModels.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(nil, plan, bpScopes, 1)
	assert.Nil(t, err)
	scopes, err := makeScopesV200(bpScopes, 1)
	assert.Nil(t, err)

	expectPlan := coreModels.PipelinePlan{
		coreModels.PipelineStage{
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

func mockBasicRes(t *testing.T) {
	jenkinsJob := &models.JenkinsJob{
		ConnectionId: 1,
		FullName:     "a/b/ccc",
	}

	scopeConfig := &models.JenkinsScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Entities: []string{"CICD"},
		},
	}

	// Refresh Global Variables and set the sql mock
	mockRes := unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.AnythingOfType("*models.JenkinsScopeConfig"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.JenkinsScopeConfig)
			*dst = *scopeConfig
		}).Return(nil)
		mockDal.On("First", mock.AnythingOfType("*models.JenkinsJob"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.JenkinsJob)
			*dst = *jenkinsJob
		}).Return(nil)
	})
	p := mockplugin.NewPluginMeta(t)
	p.On("Name").Return("dummy").Maybe()
	Init(mockRes, p)
}
