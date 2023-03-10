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
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

const PLUGIN_NAME = "fake"
const TOKEN = "this_is_a_valid_token"

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
	helper.LocalInit()
	_ = os.Setenv("REMOTE_PLUGINS_STARTUP_PATH", filepath.Join("backend/test/remote/fakeplugin/start.sh"))
	_ = os.Setenv("ENABLE_REMOTE_PLUGINS", "true")
}

func buildPython(t *testing.T) {
	fmt.Println("Build fake plugin")
	path := filepath.Join(helper.ProjectRoot, "backend/test/remote/fakeplugin/build.sh")
	cmd := exec.Command(helper.Shell, []string{path}...)
	cmd.Dir = filepath.Dir(path)
	cmd.Env = append(cmd.Env, os.Environ()...)
	r, err := utils.RunProcess(cmd,
		&utils.RunProcessOptions{
			OnStdout: func(b []byte) {
				fmt.Println(string(b))
			},
			OnStderr: func(b []byte) {
				fmt.Println(string(b))
			},
		})
	require.NoError(t, err)
	require.NoError(t, r.GetError())
}

func connectLocalServer(t *testing.T) *helper.DevlakeClient {
	fmt.Println("Connect to server")
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        helper.UseMySQL("127.0.0.1", 3307),
		CreateServer: true,
		DropDb:       true,
		Plugins:      map[string]plugin.PluginMeta{},
	})
	client.SetTimeout(60 * time.Second)
	// Wait for plugin registration
	time.Sleep(3 * time.Second)
	return client
}

func CreateTestConnection(client *helper.DevlakeClient) *helper.Connection {
	connection := client.CreateConnection(PLUGIN_NAME,
		FakePluginConnection{
			Name:  "Test connection",
			Token: TOKEN,
		},
	)

	client.SetTimeout(1)
	return connection
}

func CreateTestScope(client *helper.DevlakeClient, connectionId uint64) any {
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
	setupEnv()
	buildPython(t)
	client := connectLocalServer(t)

	CreateTestConnection(client)

	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, TOKEN, conns[0].Token)
}

func TestCreateScope(t *testing.T) {
	setupEnv()
	buildPython(t)
	client := connectLocalServer(t)
	var connectionId uint64 = 1

	CreateTestScope(client, connectionId)

	scopes := client.ListScopes(PLUGIN_NAME, connectionId)
	require.Equal(t, 1, len(scopes))
}

func TestRunPipeline(t *testing.T) {
	setupEnv()
	buildPython(t)
	client := connectLocalServer(t)
	conn := CreateTestConnection(client)

	CreateTestScope(client, conn.ID)

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
	setupEnv()
	buildPython(t)
	client := connectLocalServer(t)
	connection := CreateTestConnection(client)
	projectName := "Test project"

	client.CreateProject(&helper.ProjectConfig{
		ProjectName: projectName,
	})

	CreateTestScope(client, connection.ID)

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
