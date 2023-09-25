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
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/plugin"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/sonarqube")
	mockMeta.On("Name").Return("sonarqube").Maybe()
	err := plugin.RegisterPlugin("sonarqube", mockMeta)
	assert.Nil(t, err)
	bs := &coreModels.BlueprintScope{
		ScopeId: "f5a50c63-2e8f-4107-9014-853f6f467757",
	}

	bpScopes := make([]*coreModels.BlueprintScope, 0)
	bpScopes = append(bpScopes, bs)
	plan := make(coreModels.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(nil, plan, bpScopes, uint64(1))
	assert.Nil(t, err)
	basicRes = NewMockBasicRes()

	scopes, err := makeScopesV200(bpScopes, uint64(1))
	assert.Nil(t, err)

	expectPlan := coreModels.PipelinePlan{
		coreModels.PipelineStage{
			{
				Plugin:   "sonarqube",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"projectKey":   "f5a50c63-2e8f-4107-9014-853f6f467757",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)

	expectScopes := make([]plugin.Scope, 0)
	sonarqubeProject := &codequality.CqProject{
		DomainEntity: domainlayer.DomainEntity{
			Id: "sonarqube:SonarqubeProject:1:f5a50c63-2e8f-4107-9014-853f6f467757",
		},
		Name: "aerosolve",
	}

	expectScopes = append(expectScopes, sonarqubeProject)
	assert.Equal(t, expectScopes, scopes)
}

// NewMockBasicRes FIXME ...
func NewMockBasicRes() *mockcontext.BasicRes {
	sonarqubeProject := &models.SonarqubeProject{
		Scope: common.Scope{
			ConnectionId: 1,
		},
		ProjectKey: "f5a50c63-2e8f-4107-9014-853f6f467757",
		Name:       "aerosolve",
	}

	mockRes := new(mockcontext.BasicRes)
	mockDal := new(mockdal.Dal)

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.SonarqubeProject)
		*dst = *sonarqubeProject
	}).Return(nil).Once()

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")

	return mockRes
}
