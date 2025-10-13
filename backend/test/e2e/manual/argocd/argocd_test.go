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

package argocd

import (
	"fmt"
	"testing"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	argocdimpl "github.com/apache/incubator-devlake/plugins/argocd/impl"
	argocdmodels "github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

const pluginName = "argocd"

func TestArgoCDPlugin(t *testing.T) {
	cfg := helper.GetTestConfig[TestConfig]()
	require.NotEmpty(t, cfg.Endpoint, "Endpoint must be provided via helper.SetTestConfig(TestConfig{...})")
	require.NotEmpty(t, cfg.Token, "Token must be provided via helper.SetTestConfig(TestConfig{...})")
	require.NotEmpty(t, cfg.Applications, "At least one application is required for the ArgoCD E2E test")

	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8092,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       false,
		TruncateDb:   true,
		Plugins: []plugin.PluginMeta{
			argocdimpl.ArgoCD{},
		},
	})

	createConnection := func() *helper.Connection {
		connection := argocdmodels.ArgocdConnection{
			BaseConnection: api.BaseConnection{
				Name: "argocd-conn",
			},
			ArgocdConn: argocdmodels.ArgocdConn{
				RestConnection: api.RestConnection{
					Endpoint:         cfg.Endpoint,
					Proxy:            "",
					RateLimitPerHour: 3000,
				},
				AccessToken: api.AccessToken{
					Token: cfg.Token,
				},
			},
		}
		client.TestConnection(pluginName, connection)
		return client.CreateConnection(pluginName, connection)
	}

	connection := createConnection()

	scopeConfig := helper.Cast[argocdmodels.ArgocdScopeConfig](client.CreateScopeConfig(pluginName, connection.ID, argocdmodels.ArgocdScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Name: "argocd-default",
			Entities: []string{
				plugin.DOMAIN_TYPE_CICD,
			},
		},
		EnvNamePattern:    "(?i)prod(.*)",
		DeploymentPattern: ".*",
		ProductionPattern: "(?i)(prod|production)",
	}))

	var scopePayload []any
	for _, spec := range cfg.Applications {
		require.NotEmpty(t, spec.Name, "application name must not be empty")

		remote := client.RemoteScopes(helper.RemoteScopesQuery{
			PluginName:   pluginName,
			ConnectionId: connection.ID,
			GroupId:      spec.Project,
		})

		var selected *argocdmodels.ArgocdApplication
		for _, child := range remote.Children {
			if child.Type != "scope" {
				continue
			}
			app := helper.Cast[argocdmodels.ArgocdApplication](child.Data)
			if app.Name == spec.Name && (spec.Project == "" || app.Project == spec.Project) {
				app.ScopeConfigId = scopeConfig.ID
				selected = &app
				scopePayload = append(scopePayload, app)
				break
			}
		}
		require.NotNil(t, selected, "application %s/%s not found via remote scopes", spec.Project, spec.Name)
	}

	createdScopes := helper.Cast[[]*argocdmodels.ArgocdApplication](client.CreateScopes(pluginName, connection.ID, scopePayload...))
	require.Equal(t, len(scopePayload), len(createdScopes))

	project := client.CreateProject(&helper.ProjectConfig{
		ProjectName: fmt.Sprintf("project-%s", pluginName),
		EnableDora:  true,
	})

	listedScopes := client.ListScopes(pluginName, connection.ID, false).Scopes
	require.GreaterOrEqual(t, len(listedScopes), len(createdScopes))

	var blueprintScopes []*models.BlueprintScope
	for _, scope := range createdScopes {
		blueprintScopes = append(blueprintScopes, &models.BlueprintScope{
			ScopeId: scope.ScopeId(),
		})
	}
	require.Equal(t, len(scopePayload), len(blueprintScopes))

	require.NotNil(t, project.Blueprint, "project should have an auto-generated blueprint")
	blueprint := client.PatchBasicBlueprintV2(project.Blueprint.ID, connection.Name, &helper.BlueprintV2Config{
		Connection: &models.BlueprintConnection{
			PluginName:   pluginName,
			ConnectionId: connection.ID,
			Scopes:       blueprintScopes,
		},
		SkipOnFail:  true,
		ProjectName: project.Name,
	})
	project = client.GetProject(project.Name)
	require.Equal(t, blueprint.ID, project.Blueprint.ID)

	pipeline := client.TriggerBlueprint(project.Blueprint.ID)
	require.Equal(t, models.TASK_COMPLETED, pipeline.Status)

	client.DeleteBlueprint(project.Blueprint.ID)
	for _, scope := range createdScopes {
		client.DeleteScope(pluginName, connection.ID, scope.ScopeId(), false)
	}
	client.DeleteScopeConfig(pluginName, connection.ID, scopeConfig.ID)
	client.DeleteConnection(pluginName, connection.ID)
}
