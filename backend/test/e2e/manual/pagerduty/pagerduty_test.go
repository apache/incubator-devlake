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

package pagerduty

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/impl"
	pluginmodels "github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

const pluginName = "pagerduty"

func TestPagerDutyPlugin(t *testing.T) {
	cfg := helper.GetTestConfig[TestConfig]()
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       false,
		TruncateDb:   true,
		Plugins: []plugin.PluginMeta{
			&impl.PagerDuty{},
		},
	})
	client.SetTimeout(0)
	client.SetPipelineTimeout(0)
	connection := createConnection(cfg, client)
	t.Run("blueprint v200", func(t *testing.T) {
		serviceScopes := helper.RemoteScopesOutput{}
		var scopeData []any
		for {
			serviceScopes = client.RemoteScopes(helper.RemoteScopesQuery{
				PluginName:   pluginName,
				ConnectionId: connection.ID,
				PageToken:    serviceScopes.NextPageToken,
				Params:       nil,
			})
			for _, remoteScope := range serviceScopes.Children {
				if remoteScope.Type == "scope" {
					data := helper.Cast[pluginmodels.Service](remoteScope.Data)
					for _, serviceName := range cfg.Services {
						if serviceName == data.Name {
							scopeData = append(scopeData, &data)
						}
					}
				}
			}
			if serviceScopes.NextPageToken == "" {
				break
			}
		}
		createdScopes := helper.Cast[[]*pluginmodels.Service](client.CreateScopes(pluginName, connection.ID, scopeData...))
		require.True(t, len(createdScopes) > 0)
		uniqueString := time.Now().Format(time.RFC3339)
		outputProject := createProject(client, fmt.Sprintf("pagerduty-project-%s-%s", pluginName, uniqueString))
		var bpScopes []*plugin.BlueprintScopeV200
		for _, scope := range createdScopes {
			bpScopes = append(bpScopes, &plugin.BlueprintScopeV200{
				Id:   scope.Id,
				Name: fmt.Sprintf("pagerduty-blueprint-v200-%s", uniqueString),
			})
		}
		bp := client.CreateBasicBlueprintV2(connection.Name, &helper.BlueprintV2Config{
			Connection: &plugin.BlueprintConnectionV200{
				Plugin:       pluginName,
				ConnectionId: connection.ID,
				Scopes:       bpScopes,
			},
			TimeAfter:   cfg.TimeSince,
			SkipOnFail:  false,
			ProjectName: outputProject.Name,
		})
		// get the project ... should have a reference to the blueprint now
		outputProject = client.GetProject(outputProject.Name)
		require.Equal(t, bp.Name, outputProject.Blueprint.Name)
		fmt.Printf("=========================Triggering blueprint for project %s =========================\n", outputProject.Name)
		pipeline := client.TriggerBlueprint(bp.ID)
		require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
		createdScopesList := client.ListScopes(pluginName, connection.ID, true)
		require.True(t, len(createdScopesList) > 0)
		client.SetExpectedStatusCode(http.StatusConflict).DeleteConnection(pluginName, connection.ID)
		client.DeleteBlueprint(bp.ID)
		for _, scope := range createdScopesList {
			scopeCast := helper.Cast[pluginmodels.Service](scope.Scope)
			fmt.Printf("Deleting scope %s\n", scopeCast.Id)
			client.DeleteScope(pluginName, connection.ID, scopeCast.Id, false)
			fmt.Printf("Deleted scope %s\n", scopeCast.Id)
		}
		client.DeleteConnection(pluginName, connection.ID)
	})
	fmt.Println("======DONE======")
}

func createConnection(cfg TestConfig, client *helper.DevlakeClient) *helper.Connection {
	conn := pluginmodels.PagerDutyConn{
		RestConnection: api.RestConnection{
			Endpoint:         "https://api.pagerduty.com",
			Proxy:            "",
			RateLimitPerHour: 0,
		},
		PagerDutyAccessToken: pluginmodels.PagerDutyAccessToken{
			Token: cfg.Token,
		},
	}
	client.TestConnection(pluginName, conn)
	connections := client.ListConnections(pluginName)
	for _, connection := range connections {
		if connection.Name == "pagerduty-conn" {
			return connection
		}
	}
	return client.CreateConnection(pluginName, pluginmodels.PagerDutyConnection{
		BaseConnection: api.BaseConnection{
			Name: "pagerduty-conn",
		},
		PagerDutyConn: conn,
	})
}

func createProject(client *helper.DevlakeClient, projectName string) models.ApiOutputProject {
	projects := client.ListProjects()
	for _, project := range projects.Projects {
		if project.Name == projectName {
			outputProject := client.GetProject(projectName)
			return outputProject
		}
	}
	outputProject := client.CreateProject(&helper.ProjectConfig{
		ProjectName: projectName,
		EnableDora:  true,
	})
	return outputProject
}

func getTime(timeString string) *time.Time {
	if timeString == "" {
		return &time.Time{}
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", timeString)
	if err != nil {
		panic(err)
	}
	return &t
}
