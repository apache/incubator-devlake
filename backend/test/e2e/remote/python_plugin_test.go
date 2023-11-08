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

package remote

import (
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

func TestCreateConnection(t *testing.T) {
	client := CreateClient(t)

	CreateTestConnection(client)

	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, TOKEN, conns[0].Token)
}

func TestDeleteConnection(t *testing.T) {
	client := CreateClient(t)

	CreateTestConnection(client)

	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, TOKEN, conns[0].Token)
	refs := client.DeleteConnection(PLUGIN_NAME, conns[0].ID)
	require.Equal(t, 0, len(refs.Projects))
	require.Equal(t, 0, len(refs.Blueprints))
}

func TestDeleteConnection_Conflict(t *testing.T) {
	client := CreateClient(t)
	params := CreateTestBlueprints(t, client, 1)
	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, TOKEN, conns[0].Token)
	refs := client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(PLUGIN_NAME, params.connection.ID)
	require.Equal(t, 1, len(refs.Projects))
	require.Equal(t, 1, len(refs.Blueprints))
	client.DeleteBlueprint(params.blueprints[0].ID)
	client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(PLUGIN_NAME, params.connection.ID)
	client.DeleteScope(PLUGIN_NAME, params.connection.ID, params.scope.Id, false)
	client.DeleteConnection(PLUGIN_NAME, conns[0].ID)
}

func TestDeleteConnection_WithDependentScopesAndConfig(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	config := CreateTestScopeConfig(client, connection.ID)
	scope := CreateTestScope(client, config, connection.ID)
	refs := client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(PLUGIN_NAME, connection.ID)
	require.Equal(t, 0, len(refs.Projects))
	require.Equal(t, 0, len(refs.Blueprints))
	client.DeleteScope(PLUGIN_NAME, connection.ID, scope.Id, false)
	refs = client.DeleteConnection(PLUGIN_NAME, connection.ID)
	require.Equal(t, 0, len(refs.Projects))
	require.Equal(t, 0, len(refs.Blueprints))
	scopeRes := client.SetExpectedStatusCode(http.StatusBadRequest).ListScopes(PLUGIN_NAME, connection.ID, false).Scopes
	require.Equal(t, 0, len(scopeRes))
	configs := client.ListScopeConfigs(PLUGIN_NAME, connection.ID)
	require.Equal(t, 0, len(configs))
}

func TestRemoteScopeGroups(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)

	output := client.RemoteScopes(helper.RemoteScopesQuery{
		PluginName:   PLUGIN_NAME,
		ConnectionId: connection.ID,
	})

	scopeGroups := output.Children
	require.Equal(t, 1, len(scopeGroups))
	scope := scopeGroups[0]
	require.Equal(t, "Group 1", scope.Name)
	require.Equal(t, "group1", scope.Id)
	require.Equal(t, "group", scope.Type)
	require.Nil(t, scope.ParentId)
	require.Nil(t, scope.Data)
}

func TestRemoteScopes(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	output := client.RemoteScopes(helper.RemoteScopesQuery{
		PluginName:   PLUGIN_NAME,
		ConnectionId: connection.ID,
		GroupId:      "group1",
	})
	scopes := output.Children
	require.Equal(t, 1, len(scopes))
	scope := scopes[0]
	require.Equal(t, "Project 1", scope.Name)
	require.Equal(t, "p1", scope.Id)
	require.Equal(t, "group1", *scope.ParentId)
	require.Equal(t, "scope", scope.Type)
	require.NotNil(t, scope.Data)
	cicdScope := helper.Cast[FakeProject](scope.Data)
	require.Equal(t, connection.ID, cicdScope.ConnectionId)
	require.Equal(t, "p1", cicdScope.Id)
	require.Equal(t, "Project 1", cicdScope.Name)
	require.Equal(t, "http://fake.org/api/project/p1", cicdScope.Url)
}

func TestCreateScope(t *testing.T) {
	client := CreateClient(t)
	conn := CreateTestConnection(client)
	scopeConfig := CreateTestScopeConfig(client, conn.ID)
	scope := CreateTestScope(client, scopeConfig, conn.ID)
	scopes := client.ListScopes(PLUGIN_NAME, conn.ID, false).Scopes
	require.Equal(t, 1, len(scopes))
	cicdScope := helper.Cast[FakeProject](client.GetScope(PLUGIN_NAME, conn.ID, scope.Id, false).Scope)
	require.Equal(t, scope.Id, cicdScope.Id)
	cicdScope0 := helper.Cast[FakeProject](scopes[0].Scope)
	require.Equal(t, scope.Id, cicdScope0.Id)
	require.Equal(t, conn.ID, cicdScope.ConnectionId)
	require.Equal(t, "p1", cicdScope.Id)
	require.Equal(t, "Project 1", cicdScope.Name)
	require.Equal(t, "http://fake.org/api/project/p1", cicdScope.Url)
	cicdScope.Name = "scope-name-2"
	client.UpdateScope(PLUGIN_NAME, conn.ID, cicdScope.Id, cicdScope)
	cicdScope = helper.Cast[FakeProject](client.GetScope(PLUGIN_NAME, conn.ID, scope.Id, false).Scope)
	require.Equal(t, "scope-name-2", cicdScope.Name)
}

func TestRunPipeline(t *testing.T) {
	client := CreateClient(t)
	conn := CreateTestConnection(client)
	scopeConfig := CreateTestScopeConfig(client, conn.ID)
	scope := CreateTestScope(client, scopeConfig, conn.ID)
	pipeline := client.RunPipeline(models.NewPipeline{
		Name: "remote_test",
		Plan: []models.PipelineStage{
			{
				{
					Plugin:   PLUGIN_NAME,
					Subtasks: nil,
					Options: map[string]interface{}{
						"connectionId": conn.ID,
						"scopeId":      scope.Id,
					},
				},
			},
		},
	})
	require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
	require.Equal(t, 1, pipeline.FinishedTasks)
	require.Equal(t, "", pipeline.ErrorName)
}

func TestBlueprintV200_withScopeDeletion_Conflict(t *testing.T) {
	client := CreateClient(t)
	params := CreateTestBlueprints(t, client, 1)
	client.TriggerBlueprint(params.blueprints[0].ID)
	scopesResponse := client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	require.Equal(t, 1, len(scopesResponse[0].Blueprints))
	refs := DeleteScopeWithDataIntegrityValidation(t, client.SetExpectedStatusCode(http.StatusConflict), params.connection.ID, params.scope.Id, false)
	require.Equal(t, 1, len(refs.Blueprints))
	require.Equal(t, 1, len(refs.Projects))
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	bpsResult := client.ListBlueprints()
	require.Equal(t, 1, len(bpsResult.Blueprints))
}

func TestBlueprintV200_withBlueprintDeletion(t *testing.T) {
	client := CreateClient(t)
	params := CreateTestBlueprints(t, client, 2)
	client.TriggerBlueprint(params.blueprints[0].ID)
	scopesResponse := client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	require.Equal(t, 2, len(scopesResponse[0].Blueprints))
	client.DeleteBlueprint(params.blueprints[0].ID)
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	bpsList := client.ListBlueprints()
	require.Equal(t, 1, len(bpsList.Blueprints))
	require.Equal(t, params.blueprints[1].ID, bpsList.Blueprints[0].ID)
	projectsResponse := client.ListProjects()
	require.Equal(t, 2, len(projectsResponse.Projects))
}

func TestBlueprintV200_withBlueprintDeletion_thenScopeDeletion(t *testing.T) {
	client := CreateClient(t)
	params := CreateTestBlueprints(t, client, 1)
	client.TriggerBlueprint(params.blueprints[0].ID)
	scopesResponse := client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	require.Equal(t, 1, len(scopesResponse[0].Blueprints))
	client.DeleteBlueprint(params.blueprints[0].ID)
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	bpsList := client.ListBlueprints()
	require.Equal(t, 0, len(bpsList.Blueprints))
	refs := DeleteScopeWithDataIntegrityValidation(t, client, params.connection.ID, params.scope.Id, false)
	require.Equal(t, 0, len(refs.Blueprints))
	require.Equal(t, 0, len(refs.Projects))
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 0, len(scopesResponse))
	projectsResponse := client.ListProjects()
	require.Equal(t, 1, len(projectsResponse.Projects))
}

func TestBlueprintV200_withProjectDeletion_thenScopeDeletion(t *testing.T) {
	client := CreateClient(t)
	params := CreateTestBlueprints(t, client, 1)
	client.TriggerBlueprint(params.blueprints[0].ID)
	scopesResponse := client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	require.Equal(t, 1, len(scopesResponse[0].Blueprints))
	client.DeleteProject(params.projects[0].Name)
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 1, len(scopesResponse))
	bpsList := client.ListBlueprints()
	require.Equal(t, 0, len(bpsList.Blueprints))
	refs := DeleteScopeWithDataIntegrityValidation(t, client, params.connection.ID, params.scope.Id, false)
	require.Equal(t, 0, len(refs.Blueprints))
	require.Equal(t, 0, len(refs.Projects))
	scopesResponse = client.ListScopes(PLUGIN_NAME, params.connection.ID, true).Scopes
	require.Equal(t, 0, len(scopesResponse))
	projectsResponse := client.ListProjects()
	require.Equal(t, 0, len(projectsResponse.Projects))
}

func TestCreateScopeConfig(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	scopeConfig := FakeScopeConfig{Name: "Scope config", Env: "test env", Entities: []string{plugin.DOMAIN_TYPE_CICD}}

	res := client.CreateScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig)
	scopeConfig = helper.Cast[FakeScopeConfig](res)

	res = client.GetScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig.Id)
	scopeConfig = helper.Cast[FakeScopeConfig](res)
	require.Equal(t, "Scope config", scopeConfig.Name)
	require.Equal(t, "test env", scopeConfig.Env)
	require.Equal(t, []string{plugin.DOMAIN_TYPE_CICD}, scopeConfig.Entities)
}

func TestUpdateScopeConfig(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	res := client.CreateScopeConfig(PLUGIN_NAME, connection.ID, FakeScopeConfig{Name: "old name", Env: "old env", Entities: []string{}})
	oldScopeConfig := helper.Cast[FakeScopeConfig](res)

	client.PatchScopeConfig(PLUGIN_NAME, connection.ID, oldScopeConfig.Id, FakeScopeConfig{Name: "new name", Env: "new env", Entities: []string{plugin.DOMAIN_TYPE_CICD}})

	res = client.GetScopeConfig(PLUGIN_NAME, connection.ID, oldScopeConfig.Id)
	scopeConfig := helper.Cast[FakeScopeConfig](res)
	require.Equal(t, "new name", scopeConfig.Name)
	require.Equal(t, "new env", scopeConfig.Env)
}

func TestDeleteScopeConfig(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	scopeConfig := FakeScopeConfig{Name: "Scope config", Env: "test env", Entities: []string{plugin.DOMAIN_TYPE_CICD}}
	scopeConfig = helper.Cast[FakeScopeConfig](client.CreateScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig))

	configs := helper.Cast[[]FakeScopeConfig](client.ListScopeConfigs(PLUGIN_NAME, connection.ID))
	require.Equal(t, 1, len(configs))

	client.DeleteScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig.Id)
	configs = helper.Cast[[]FakeScopeConfig](client.ListScopeConfigs(PLUGIN_NAME, connection.ID))
	require.Equal(t, 0, len(configs))
}

func TestDeleteScopeConfig_WithReferencingScope(t *testing.T) {
	client := CreateClient(t)
	connection := CreateTestConnection(client)
	scopeConfig := FakeScopeConfig{Name: "Scope config", Env: "test env", Entities: []string{plugin.DOMAIN_TYPE_CICD}}
	scopeConfig = helper.Cast[FakeScopeConfig](client.CreateScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig))

	scope := CreateTestScope(client, &scopeConfig, connection.ID)
	require.Equal(t, scopeConfig.Id, scope.ScopeConfigId)

	client.DeleteScopeConfig(PLUGIN_NAME, connection.ID, scopeConfig.Id)
	scope = helper.Cast[*FakeProject](client.GetScope(PLUGIN_NAME, connection.ID, scope.Id, false))
	require.Equal(t, uint64(0), scope.ScopeConfigId)

}
