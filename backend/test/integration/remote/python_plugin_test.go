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

package test

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/test/integration/helper"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/stretchr/testify/require"
)

const (
	PLUGIN_NAME     = "fake"
	TOKEN           = "this_is_a_valid_token"
	FAKE_PLUGIN_DIR = "python/test/fakeplugin"
)

type FakePluginConnection struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type FakeProject struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	ConnectionId         uint64 `json:"connection_id"`
	TransformationRuleId uint64 `json:"transformation_rule_id"`
}

type FakeTxRule struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
}

func setupEnv() {
	fmt.Println("Setup test env")
	helper.Init()
	path := filepath.Join(helper.ProjectRoot, FAKE_PLUGIN_DIR, "start.sh") // make sure the path is correct
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("REMOTE_PLUGINS_STARTUP_PATH", path)
	_ = os.Setenv("ENABLE_REMOTE_PLUGINS", "true")
}

func connectLocalServer(t *testing.T) *helper.DevlakeClient {
	fmt.Println("Connect to server")
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		TruncateDb:   true,
		Plugins:      map[string]plugin.PluginMeta{},
	})
	client.SetTimeout(60 * time.Second)
	// Wait for plugin registration
	time.Sleep(3 * time.Second)
	return client
}

func createClient(t *testing.T) *helper.DevlakeClient {
	setupEnv()
	return connectLocalServer(t)
}

func createTestConnection(client *helper.DevlakeClient) *helper.Connection {
	connection := client.CreateConnection(PLUGIN_NAME,
		FakePluginConnection{
			Name:  "Test connection",
			Token: TOKEN,
		},
	)

	client.SetTimeout(1)
	return connection
}

func createTestScope(client *helper.DevlakeClient, connectionId uint64) any {
	res := client.CreateTransformationRule(PLUGIN_NAME, FakeTxRule{Name: "Tx rule", Env: "test env"})
	rule, ok := res.(map[string]interface{})
	if !ok {
		panic("Cannot cast transform rule")
	}
	ruleId := uint64(rule["id"].(float64))

	scope := client.CreateScope(PLUGIN_NAME,
		connectionId,
		FakeProject{
			Id:                   "12345",
			Name:                 "Test project",
			ConnectionId:         connectionId,
			TransformationRuleId: ruleId,
		},
	)

	client.SetTimeout(1)
	return scope
}

func TestCreateConnection(t *testing.T) {
	client := createClient(t)

	createTestConnection(client)

	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, TOKEN, conns[0].Token)
}

func TestRemoteScopeGroups(t *testing.T) {
	client := createClient(t)
	connection := createTestConnection(client)

	output := client.RemoteScopes(helper.RemoteScopesQuery{
		PluginName:   PLUGIN_NAME,
		ConnectionId: connection.ID,
	})

	scopeGroups := output.Children
	require.Equal(t, 1, len(scopeGroups))
	scope := scopeGroups[0]
	require.Equal(t, "Group 1", scope.Name)
	require.Equal(t, "group1", scope.Id)
}

func TestRemoteScopes(t *testing.T) {
	client := createClient(t)
	connection := createTestConnection(client)

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
}

func TestCreateScope(t *testing.T) {
	client := createClient(t)
	var connectionId uint64 = 1

	createTestScope(client, connectionId)

	scopes := client.ListScopes(PLUGIN_NAME, connectionId)
	require.Equal(t, 1, len(scopes))
}

func TestRunPipeline(t *testing.T) {
	t.SkipNow() //Fix later
	client := createClient(t)
	conn := createTestConnection(client)

	createTestScope(client, conn.ID)

	pipeline := client.RunPipeline(models.NewPipeline{
		Name: "remote_test",
		Plan: []plugin.PipelineStage{
			{
				{
					Plugin:   PLUGIN_NAME,
					Subtasks: nil,
					Options: map[string]interface{}{
						"connectionId": conn.ID,
						"scopeId":      "12345",
					},
				},
			},
		},
	})

	require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
	require.Equal(t, 1, pipeline.FinishedTasks)
	require.Equal(t, "", pipeline.ErrorName)
}

func TestBlueprintV200(t *testing.T) {
	t.SkipNow() //Fix later
	client := createClient(t)
	connection := createTestConnection(client)
	projectName := "Test project"
	client.CreateProject(&helper.ProjectConfig{
		ProjectName: projectName,
	})
	createTestScope(client, connection.ID)

	blueprint := client.CreateBasicBlueprintV2(
		"Test blueprint",
		&helper.BlueprintV2Config{
			Connection: &plugin.BlueprintConnectionV200{
				Plugin:       "fake",
				ConnectionId: connection.ID,
				Scopes: []*plugin.BlueprintScopeV200{
					{
						Id:   "12345",
						Name: "Test scope",
						Entities: []string{
							plugin.DOMAIN_TYPE_CICD,
						},
					},
				},
			},
			SkipOnFail:  true,
			ProjectName: projectName,
		},
	)

	project := client.GetProject(projectName)
	require.Equal(t, blueprint.Name, project.Blueprint.Name)
	client.TriggerBlueprint(blueprint.ID)
}
