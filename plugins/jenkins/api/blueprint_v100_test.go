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
	"encoding/json"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessScope(t *testing.T) {
	connection := &models.JenkinsConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Name: "jenkins-test",
				Model: common.Model{
					ID: 1,
				},
			},
			Endpoint:         "https://api.github.com/",
			Proxy:            "",
			RateLimitPerHour: 0,
		},
		BasicAuth: helper.BasicAuth{
			Username: "Username",
			Password: "Password",
		},
	}

	bs := &core.BlueprintScopeV100{
		Entities: []string{"CICD"},
		Options: json.RawMessage(`{
              "jobName": "testJob"
            }`),
		Transformation: json.RawMessage(`{
              "productionPattern": "(?i)build-and-deploy",
              "deploymentPattern": "deploy"
            }`),
	}
	scopes := make([]*core.BlueprintScopeV100, 0)
	scopes = append(scopes, bs)
	plan, err := makePipelinePlanV100(nil, scopes, connection)
	assert.Nil(t, err)

	expectPlan := core.PipelinePlan{
		core.PipelineStage{
			{
				Plugin:   "jenkins",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"jobName":      "testJob",
					"connectionId": uint64(1),
					"transformationRules": map[string]interface{}{
						"deploymentPattern": "deploy",
					},
				},
			},
		},
		core.PipelineStage{
			{
				Plugin:   "dora",
				Subtasks: []string{"EnrichTaskEnv"},
				Options: map[string]interface{}{
					"prefix": "jenkins",
					"transformationRules": map[string]interface{}{
						"productionPattern": "(?i)build-and-deploy",
					},
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)
}
