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

package testmo

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/impl"
	pluginmodels "github.com/apache/incubator-devlake/plugins/testmo/models"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

const pluginName = "testmo"

func TestTestmoPlugin(t *testing.T) {
	cfg := helper.GetTestConfig[TestConfig]()

	// Disable remote plugins to avoid poetry dependency
	config.GetConfig().Set("DISABLED_REMOTE_PLUGINS", true)

	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       true, // Use fresh database for each test run
		TruncateDb:   true,
		Plugins: []plugin.PluginMeta{
			&impl.Testmo{},
		},
	})
	client.SetTimeout(0)
	client.SetPipelineTimeout(0)
	connection := createConnection(cfg, client)
	t.Run("blueprint v200", func(t *testing.T) {
		projectScopes := helper.RemoteScopesOutput{}
		var scopeData []any
		for {
			projectScopes = client.RemoteScopes(helper.RemoteScopesQuery{
				PluginName:   pluginName,
				ConnectionId: connection.ID,
				PageToken:    projectScopes.NextPageToken,
				Params:       nil,
			})
			for _, remoteScope := range projectScopes.Children {
				if remoteScope.Type == "scope" {
					data := helper.Cast[pluginmodels.TestmoProject](remoteScope.Data)
					for _, projectName := range cfg.Projects {
						if projectName == data.Name {
							scopeData = append(scopeData, &data)
						}
					}
				}
			}
			if projectScopes.NextPageToken == "" {
				break
			}
		}
		createdScopes := helper.Cast[[]*pluginmodels.TestmoProject](client.CreateScopes(pluginName, connection.ID, scopeData...))
		require.True(t, len(createdScopes) > 0)

		projectName := fmt.Sprintf("testmo-project-%s-%s", pluginName, "unique")

		outputProject := client.CreateProject(&helper.ProjectConfig{
			ProjectName: projectName,
			EnableDora:  true,
		})

		// Delete any existing blueprints for this project
		if outputProject.Blueprint != nil {
			client.DeleteBlueprint(outputProject.Blueprint.ID)
		}

		// Create our blueprint
		var bpScopes []*models.BlueprintScope
		for _, scope := range createdScopes {
			bpScopes = append(bpScopes, &models.BlueprintScope{
				ScopeId: scope.ScopeId(),
			})
		}
		bp := client.CreateBasicBlueprintV2(connection.Name, &helper.BlueprintV2Config{
			Connection: &models.BlueprintConnection{
				PluginName:   pluginName,
				ConnectionId: connection.ID,
				Scopes:       bpScopes,
			},
			TimeAfter:   cfg.TimeSince,
			SkipOnFail:  false,
			ProjectName: outputProject.Name,
		})

		// get the project ... should have a reference to the blueprint now
		outputProject = client.GetProject(projectName)
		require.Equal(t, bp.Name, outputProject.Blueprint.Name)
		fmt.Printf("=========================Triggering blueprint for project %s =========================\n", outputProject.Name)
		pipeline := client.TriggerBlueprint(bp.ID)
		require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
		createdScopesList := client.ListScopes(pluginName, connection.ID, true).Scopes
		require.True(t, len(createdScopesList) > 0)
		client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(pluginName, connection.ID)
		client.DeleteBlueprint(bp.ID)
		for _, scope := range createdScopesList {
			// Extract ID directly from the nested scope data
			var scopeId string
			if scopeWrapper, ok := scope.Scope.(map[string]interface{}); ok {
				if scopeData, exists := scopeWrapper["scope"]; exists {
					if scopeMap, ok := scopeData.(map[string]interface{}); ok {
						if rawId, exists := scopeMap["id"]; exists {
							scopeId = fmt.Sprintf("%v", rawId)
						}
					}
				}
			}

			if scopeId != "0" && scopeId != "" {
				fmt.Printf("Deleting scope %s\n", scopeId)
				client.DeleteScope(pluginName, connection.ID, scopeId, false)
				fmt.Printf("Deleted scope %s\n", scopeId)
			} else {
				fmt.Printf("Skipping scope with invalid ID: %s\n", scopeId)
			}
		}
		client.DeleteConnection(pluginName, connection.ID)
	})
	fmt.Println("======DONE======")
}

func createConnection(cfg TestConfig, client *helper.DevlakeClient) *helper.Connection {
	conn := pluginmodels.TestmoConn{
		RestConnection: api.RestConnection{
			Endpoint:         "https://yourorganization.testmo.net/api/v1", // Update this endpoint in testmo_local_test.go as needed
			Proxy:            "",
			RateLimitPerHour: 0,
		},
		AccessToken: api.AccessToken{
			Token: cfg.Token,
		},
	}
	client.TestConnection(pluginName, conn)
	connections := client.ListConnections(pluginName)
	for _, connection := range connections {
		if connection.Name == "testmo-conn" {
			return connection
		}
	}
	return client.CreateConnection(pluginName, pluginmodels.TestmoConnection{
		BaseConnection: api.BaseConnection{
			Name: "testmo-conn",
		},
		TestmoConn: conn,
	})
}
