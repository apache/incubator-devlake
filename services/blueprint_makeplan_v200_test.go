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

package services

import (
	"encoding/json"
	"testing"

	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
)

func TestMakePlanV200(t *testing.T) {
	const projectName = "TestMakePlanV200-project"
	githubName := "TestMakePlanV200-github" // mimic github
	// mock github plugin as a data source plugin
	githubConnId := uint64(1)
	githubScopes := []*core.BlueprintScopeV200{
		{Id: "", Name: "apache/incubator-devlake"},
		{Id: "", Name: "apache/incubator-devlake-website"},
	}
	githubOutputPlan := core.PipelinePlan{
		{
			{Plugin: githubName, Options: map[string]interface{}{"name": "apache/incubator-devlake"}},
			{Plugin: "gitextractor", Options: map[string]interface{}{"url": "http://gihub.com/apache/incubator-devlake.git"}},
		},
		{
			{Plugin: githubName, Options: map[string]interface{}{"name": "apache/incubator-devlake-website"}},
			{Plugin: "gitextractor", Options: map[string]interface{}{"url": "http://gihub.com/apache/incubator-devlake-website.git"}},
		},
	}
	githubOutputScopes := []core.Scope{
		&code.Repo{DomainEntity: domainlayer.DomainEntity{Id: "github:GithubRepo:1:123"}, Name: "apache/incubator-devlake"},
		&ticket.Board{DomainEntity: domainlayer.DomainEntity{Id: "github:GithubRepo:1:123"}, Name: "apache/incubator-devlake"},
	}
	github := new(mocks.CompositeDataSourcePluginBlueprintV200)
	github.On("MakeDataSourcePipelinePlanV200", githubConnId, githubScopes).Return(githubOutputPlan, githubOutputScopes, nil)

	// mock dora plugin as a metric plugin
	doraName := "TestMakePlanV200-dora"
	doraOutputPlan := core.PipelinePlan{
		{
			{Plugin: "refdiff", Subtasks: []string{"calculateDeploymentDiffs"}, Options: map[string]interface{}{"projectName": projectName}},
			{Plugin: doraName},
		},
	}
	dora := new(mocks.CompositeMetricPluginBlueprintV200)
	dora.On("MakeMetricPluginPipelinePlanV200", projectName, json.RawMessage(nil)).Return(doraOutputPlan, nil)

	// expectation, establish expectation before any code being launch to avoid unwanted modification
	expectedPlan := make(core.PipelinePlan, 0)
	expectedPlan = append(expectedPlan, githubOutputPlan...)
	expectedPlan = append(expectedPlan, doraOutputPlan...)
	expectedScopes := append(make([]core.Scope, 0), githubOutputScopes...)

	// plugin registration
	core.RegisterPlugin(githubName, github)
	core.RegisterPlugin(doraName, dora)

	// put them together and call GeneratePlanJsonV200
	connections, _ := json.Marshal([]*core.BlueprintConnectionV200{
		{Plugin: githubName, ConnectionId: githubConnId, Scopes: githubScopes},
	})
	sources := &models.BlueprintSettings{
		Version:     "2.0.0",
		Connections: connections,
	}
	metrics := map[string]json.RawMessage{
		doraName: nil,
	}

	plan, scopes, err := GeneratePlanJsonV200(projectName, sources, metrics)
	assert.Nil(t, err)

	assert.Equal(t, expectedPlan, plan)
	assert.Equal(t, expectedScopes, scopes)
}
