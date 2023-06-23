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

package gitlab

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	gitextractor "github.com/apache/incubator-devlake/plugins/gitextractor/impl"
	gitlab "github.com/apache/incubator-devlake/plugins/gitlab/impl"
	pluginmodels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const pluginName = "gitlab"

func TestGitlabPlugin(t *testing.T) {
	createConnection := func(cfg TestConfig, client *helper.DevlakeClient) *helper.Connection {
		conn := pluginmodels.GitlabConnection{
			BaseConnection: api.BaseConnection{
				Name: "gitlab-conn",
			},
			GitlabConn: pluginmodels.GitlabConn{
				RestConnection: api.RestConnection{
					Endpoint:         "https://gitlab.com/api/v4",
					Proxy:            "",
					RateLimitPerHour: 0,
				},
				AccessToken: api.AccessToken{Token: cfg.Token},
			},
		}
		client.TestConnection(pluginName, conn)
		return client.CreateConnection(pluginName, conn)
	}
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       false,
		TruncateDb:   true,
		Plugins: map[string]plugin.PluginMeta{
			"gitlab":       gitlab.Gitlab(""),
			"gitextractor": gitextractor.GitExtractor{},
		},
	})
	cfg := helper.GetTestConfig[TestConfig]()
	connection := createConnection(cfg, client)
	t.Run("blueprint v200", func(t *testing.T) {
		scopeConfig := helper.Cast[pluginmodels.GitlabScopeConfig](client.CreateScopeConfig("gitlab", connection.ID,
			pluginmodels.GitlabScopeConfig{
				ScopeConfig: common.ScopeConfig{
					Entities: []string{
						plugin.DOMAIN_TYPE_CICD,
						plugin.DOMAIN_TYPE_CODE,
						plugin.DOMAIN_TYPE_CODE_REVIEW,
					},
				},
				Name:                 "config-1",
				PrType:               "",
				PrComponent:          "",
				PrBodyClosePattern:   "",
				IssueSeverity:        "",
				IssuePriority:        "",
				IssueComponent:       "",
				IssueTypeBug:         "",
				IssueTypeIncident:    "",
				IssueTypeRequirement: "",
				DeploymentPattern:    ".*",
				ProductionPattern:    ".*",             // this triggers dora
				Refdiff:              map[string]any{}, // this is technically a true/false (nil or not)
			}))
		_ = scopeConfig
		remoteScopes := client.RemoteScopes(helper.RemoteScopesQuery{
			PluginName:   pluginName,
			ConnectionId: connection.ID,
			PageToken:    "",
			Params:       nil,
		})
		{
			// this doesn't have any direct use-case (for testing anyway)
			searchRemoteScopes := client.SearchRemoteScopes(helper.SearchRemoteScopesQuery{
				PluginName:   pluginName,
				ConnectionId: connection.ID,
				Search:       "projects",
				Page:         10,
				PageSize:     5,
				Params:       nil,
			})
			_ = searchRemoteScopes
		}
		var scopeData []any
		for _, remoteScope := range remoteScopes.Children {
			if remoteScope.Type == "scope" {
				data := helper.Cast[pluginmodels.GitlabProject](remoteScope.Data)
				if len(cfg.Projects) == 0 || helper.Contains(cfg.Projects, data.Name) {
					data.ScopeConfigId = scopeConfig.ID
					scopeData = append(scopeData, data)
				}
			}
		}
		createdScopes := helper.Cast[[]*pluginmodels.GitlabProject](client.CreateScopes(pluginName, connection.ID, scopeData...))
		listedScopes := client.ListScopes(pluginName, connection.ID, false)
		require.Equal(t, len(createdScopes), len(listedScopes))
		outputProject := client.CreateProject(&helper.ProjectConfig{
			ProjectName: fmt.Sprintf("project-%s", pluginName),
			EnableDora:  true,
		})
		projects := client.ListProjects()
		require.Equal(t, 1, len(projects.Projects))
		var scopes []*plugin.BlueprintScopeV200
		for _, scope := range listedScopes {
			project := helper.Cast[pluginmodels.GitlabProject](scope.Scope)
			scopes = append(scopes, &plugin.BlueprintScopeV200{
				Id:   fmt.Sprintf("%s", project.ScopeId()),
				Name: "blueprint-v200",
			})
		}
		bp := client.CreateBasicBlueprintV2(connection.Name, &helper.BlueprintV2Config{
			Connection: &plugin.BlueprintConnectionV200{
				Plugin:       pluginName,
				ConnectionId: connection.ID,
				Scopes:       scopes,
			},
			SkipOnFail:  true,
			ProjectName: outputProject.Name,
		})
		// get the project ... should have a reference to the blueprint now
		outputProject = client.GetProject(outputProject.Name)
		require.Equal(t, bp.Name, outputProject.Blueprint.Name)
		fmt.Printf("=========================Triggering blueprint for project %s =========================\n", outputProject.Name)
		pipeline := client.TriggerBlueprint(bp.ID)
		require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
		client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(pluginName, connection.ID)
		client.DeleteScopeConfig(pluginName, connection.ID, scopeConfig.ID)
		client.DeleteBlueprint(bp.ID)
		client.DeleteConnection(pluginName, connection.ID)
	})
	fmt.Println("======DONE======")
}
