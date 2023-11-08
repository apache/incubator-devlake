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

package azuredevops

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	gitextractor "github.com/apache/incubator-devlake/plugins/gitextractor/impl"
	pluginmodels "github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

const (
	azurePlugin = "azuredevops"
)

func TestAzure(t *testing.T) {
	cfg := helper.GetTestConfig[TestConfig]()
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       false,
		TruncateDb:   true,
		Plugins: []plugin.PluginMeta{
			gitextractor.GitExtractor{},
		},
	})
	client.SetTimeout(60 * time.Second)
	// Wait for plugin registration
	time.Sleep(5 * time.Second)
	fmt.Println("Create new connection")
	connection := client.CreateConnection(azurePlugin,
		AzureConnection{
			Name:         "azure_conn",
			Token:        cfg.Token,
			Organization: cfg.Org,
		},
	)
	client.SetTimeout(0)
	client.SetPipelineTimeout(0)
	conns := client.ListConnections(azurePlugin)
	require.Greater(t, len(conns), 0)
	t.Run("v200-plugin", func(t *testing.T) {
		outputProject := client.CreateProject(&helper.ProjectConfig{
			ProjectName: "project-1",
			EnableDora:  true,
		})
		projects := client.ListProjects()
		require.Equal(t, 1, len(projects.Projects))
		project := projects.Projects[0]
		repoConfig := helper.Cast[AzureGitRepositoryConfig](client.CreateScopeConfig(azurePlugin, connection.ID,
			AzureGitRepositoryConfig{
				ScopeConfig: common.ScopeConfig{
					Entities: []string{
						plugin.DOMAIN_TYPE_CICD,
						plugin.DOMAIN_TYPE_CODE,
						plugin.DOMAIN_TYPE_CODE_REVIEW,
					},
				},
				Refdiff: Refdiff{
					TagsPattern: ".*",
					TagsLimit:   1,
					TagsOrder:   "",
				},
				DeploymentPattern: ".*",
				ProductionPattern: ".*",
			},
		))
		_ = repoConfig
		remoteScopes := client.RemoteScopes(helper.RemoteScopesQuery{
			PluginName:   azurePlugin,
			ConnectionId: connection.ID,
			GroupId:      fmt.Sprintf("%s/%s", cfg.Org, cfg.Project),
			PageToken:    "",
			Params:       nil,
		})
		scopes := helper.Cast[[]AzureGitRepo](client.CreateScopes(azurePlugin, connection.ID, remoteScopesToScopes(remoteScopes, cfg.Repos)...))
		scopesCount := len(scopes)
		scopesResponse := client.ListScopes(azurePlugin, connection.ID, false)
		require.Equal(t, scopesCount, len(scopesResponse.Scopes))
		// associate scopes with the scope config
		for _, scope := range scopes {
			scope.ScopeConfigId = repoConfig.ID
			scope = helper.Cast[AzureGitRepo](client.UpdateScope(azurePlugin, connection.ID, scope.Id, scope))
			require.Equal(t, repoConfig.ID, scope.ScopeConfigId)
		}
		// create bp_scopes from the scopes
		var bpScopes []*models.BlueprintScope
		for _, scope := range scopes {
			bpScopes = append(bpScopes, &models.BlueprintScope{
				ScopeId: scope.Id,
			})
		}
		// create the bp
		bp := client.CreateBasicBlueprintV2(connection.Name, &helper.BlueprintV2Config{
			Connection: &models.BlueprintConnection{
				PluginName:   azurePlugin,
				ConnectionId: connection.ID,
				Scopes:       bpScopes,
			},
			TimeAfter:   nil,
			SkipOnFail:  false,
			ProjectName: project.Name,
		})
		// get the project ... should have a reference to the blueprint now
		outputProject = client.GetProject(project.Name)
		require.Equal(t, bp.Name, outputProject.Blueprint.Name)
		// run the bp
		pipeline := client.TriggerBlueprint(bp.ID)
		require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
		createdScopesList := client.ListScopes(azurePlugin, connection.ID, true).Scopes
		require.True(t, len(createdScopesList) > 0)
		client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(azurePlugin, connection.ID)
		client.DeleteScopeConfig(azurePlugin, connection.ID, repoConfig.ID)
		client.DeleteBlueprint(bp.ID)
		for _, scope := range createdScopesList {
			scopeCast := helper.Cast[pluginmodels.Service](scope.Scope)
			fmt.Printf("Deleting scope %s\n", scopeCast.Id)
			client.DeleteScope(azurePlugin, connection.ID, scopeCast.Id, false)
			fmt.Printf("Deleted scope %s\n", scopeCast.Id)
		}
		client.DeleteConnection(azurePlugin, connection.ID)
	})
	fmt.Println("========DONE=======")
}

func remoteScopesToScopes(remoteScopes helper.RemoteScopesOutput, filters []string) []any {
	var a []any
	for _, c := range remoteScopes.Children {
		repo := helper.Cast[AzureGitRepo](c.Data)
		if len(filters) == 0 {
			a = append(a, repo)
		} else {
			for _, f := range filters {
				if len(filters) == 0 || strings.Contains(repo.Name, f) {
					a = append(a, repo)
				}
			}
		}
	}
	return a
}
