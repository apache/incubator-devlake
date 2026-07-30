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

	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/org/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakePlanV200(t *testing.T) {
	const projectName = "TestMakePlanV200-project"
	githubName := "TestMakePlanV200-github" // mimic github
	// mock github plugin as a data source plugin
	githubConnId := uint64(1)
	githubScopes := []*coreModels.BlueprintScope{
		{ScopeId: "github:GithubRepo:1:123"},
		{ScopeId: "github:GithubRepo:1:321"},
	}
	githubOutputPlan := coreModels.PipelinePlan{
		{
			{Plugin: githubName, Options: map[string]interface{}{"name": "apache/incubator-devlake"}},
			{Plugin: "gitextractor", Options: map[string]interface{}{"url": "http://gihub.com/apache/incubator-devlake.git"}},
		},
		{
			{Plugin: githubName, Options: map[string]interface{}{"name": "apache/incubator-devlake-website"}},
			{Plugin: "gitextractor", Options: map[string]interface{}{"url": "http://gihub.com/apache/incubator-devlake-website.git"}},
		},
	}
	githubOutputScopes := []plugin.Scope{
		&code.Repo{DomainEntity: domainlayer.DomainEntity{Id: "github:GithubRepo:1:123"}, Name: "apache/incubator-devlake"},
		&ticket.Board{DomainEntity: domainlayer.DomainEntity{Id: "github:GithubRepo:1:123"}, Name: "apache/incubator-devlake"},
	}
	github := new(mockplugin.CompositeDataSourcePluginBlueprintV200)
	github.On("MakeDataSourcePipelinePlanV200", githubConnId, githubScopes).Return(githubOutputPlan, githubOutputScopes, nil)

	// mock dora plugin as a metric plugin
	doraName := "TestMakePlanV200-dora"
	doraOutputPlan := coreModels.PipelinePlan{
		{
			{Plugin: "refdiff", Subtasks: []string{"calculateProjectDeploymentCommitsDiff"}, Options: map[string]interface{}{"projectName": projectName}},
			{Plugin: doraName},
		},
	}
	dora := new(mockplugin.CompositeMetricPluginBlueprintV200)
	dora.On("MakeMetricPluginPipelinePlanV200", projectName, json.RawMessage("{}")).Return(doraOutputPlan, nil)

	// mock org plugin
	org := new(mockplugin.CompositeProjectMapper)
	orgPlan := coreModels.PipelinePlan{
		{
			{Plugin: "org", Subtasks: []string{"setProjectMapping"}, Options: map[string]interface{}{"projectMappings": []interface{}{tasks.NewProjectMapping(projectName, githubOutputScopes)}}},
		},
	}
	org.On("MapProject", projectName, githubOutputScopes).Return(orgPlan, nil)

	// expectation, establish expectation before any code being launch to avoid unwanted modification
	expectedPlan := make(coreModels.PipelinePlan, 0)
	expectedPlan = append(expectedPlan, orgPlan...)
	expectedPlan = append(expectedPlan, githubOutputPlan...)
	expectedPlan = append(expectedPlan, doraOutputPlan...)

	// plugin registration
	plugin.RegisterPlugin(githubName, github)
	plugin.RegisterPlugin(doraName, dora)
	plugin.RegisterPlugin("org", org)

	// put them together and call GeneratePlanJsonV200
	connections := []*coreModels.BlueprintConnection{
		{PluginName: githubName, ConnectionId: githubConnId, Scopes: githubScopes},
	}
	metrics := map[string]json.RawMessage{
		doraName: nil,
	}

	plan, err := GeneratePlanJsonV200(projectName, connections, metrics, false)
	assert.Nil(t, err)

	assert.Equal(t, expectedPlan, plan)
}

// TestMakePlanV200InjectsProjectNameIntoDataSourceTaskOptions verifies that
// GeneratePlanJsonV200 stamps the running Devlake project's name into every
// data-source task's options, without requiring the plugin's own
// MakeDataSourcePipelinePlanV200 implementation to know about it. This is
// what lets a scope shared by multiple projects (e.g. a Jira board added to
// two different Devlake projects) resolve the correct project per pipeline
// run instead of guessing from a many-to-one reverse lookup.
func TestMakePlanV200InjectsProjectNameIntoDataSourceTaskOptions(t *testing.T) {
	const projectName = "TestMakePlanV200InjectsProjectName-project"
	jiraName := "TestMakePlanV200InjectsProjectName-jira"
	connId := uint64(1)
	scopes := []*coreModels.BlueprintScope{{ScopeId: "jira:JiraBoard:1:1"}}
	outputPlan := coreModels.PipelinePlan{
		{
			{Plugin: jiraName, Options: map[string]interface{}{"boardId": float64(1)}},
		},
	}
	jira := new(mockplugin.CompositeDataSourcePluginBlueprintV200)
	jira.On("MakeDataSourcePipelinePlanV200", connId, scopes).Return(outputPlan, []plugin.Scope(nil), nil)
	plugin.RegisterPlugin(jiraName, jira)

	// GeneratePlanJsonV200 also calls the "org" plugin's ProjectMapper
	// whenever projectName != "". Re-register it here (overwriting any
	// registration from other tests in this package) with permissive
	// matchers, since this test doesn't care about project_mapping.
	org := new(mockplugin.CompositeProjectMapper)
	org.On("MapProject", mock.Anything, mock.Anything).Return(coreModels.PipelinePlan{}, nil)
	plugin.RegisterPlugin("org", org)

	connections := []*coreModels.BlueprintConnection{
		{PluginName: jiraName, ConnectionId: connId, Scopes: scopes},
	}

	plan, err := GeneratePlanJsonV200(projectName, connections, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(plan))
	assert.Equal(t, 1, len(plan[0]))
	assert.Equal(t, projectName, plan[0][0].Options["projectName"])
	// the plugin's own option is preserved alongside the injected one
	assert.Equal(t, float64(1), plan[0][0].Options["boardId"])
}
