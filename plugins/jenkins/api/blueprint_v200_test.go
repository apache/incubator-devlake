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
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/jenkins")
	err := core.RegisterPlugin("jenkins", mockMeta)
	assert.Nil(t, err)
	bs := &core.BlueprintScopeV200{
		Entities: []string{"CICD"},
		Id:       "",
		Name:     "",
	}
	bpScopes := make([]*core.BlueprintScopeV200, 0)
	bpScopes = append(bpScopes, bs)

	plan := make(core.PipelinePlan, len(bpScopes))
	scopes := make([]core.Scope, 0, len(bpScopes))
	for i, bpScope := range bpScopes {
		jenkinsJob := &models.JenkinsJob{
			ConnectionId: 1,
			FullName:     "a/b/ccc",
			Path:         "job/a/job/b/",
			Name:         "ccc",
		}

		transformationRule := &models.JenkinsTransformationRule{
			Model: common.Model{
				ID: 1,
			},
			Name:              "jenkins transformation rule",
			DeploymentPattern: "hey,man,wasup",
		}

		var scope []core.Scope
		plan[i], scope, err = makeDataSourcePipelinePlanV200(nil, bpScope, jenkinsJob, transformationRule)

		assert.Nil(t, err)
		scopes = append(scopes, scope...)
	}
	expectPlan := core.PipelinePlan{
		core.PipelineStage{
			{
				Plugin:     "jenkins",
				Subtasks:   []string{},
				SkipOnFail: false,
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"jobFullName":  "a/b/ccc",
					"transformationRules": map[string]interface{}{
						"name":              "jenkins transformation rule",
						"deploymentPattern": "hey,man,wasup",
					},
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)

	expectScopes := make([]core.Scope, 0)
	scopeCicd := &devops.CicdScope{
		DomainEntity: domainlayer.DomainEntity{
			Id: "jenkins:JenkinsJob:1:a/b/ccc",
		},
		Name:        "a/b/ccc",
		Description: "",
		Url:         "",
		CreatedDate: nil,
	}

	expectScopes = append(expectScopes, scopeCicd)
	assert.Equal(t, expectScopes, scopes)
}
