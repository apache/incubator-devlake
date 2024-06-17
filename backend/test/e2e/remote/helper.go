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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/services"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/stretchr/testify/require"

	"github.com/apache/incubator-devlake/test/helper"
)

const (
	PLUGIN_NAME     = "fake"
	TOKEN           = "this_is_a_valid_token"
	EMPTY_TOKEN     = ""
	FAKE_PLUGIN_DIR = "python/test/fakeplugin"
)

var (
	PluginDataTables = []string{
		"_raw_fake_fakepipelinestream",
		"_tool_fakeplugin_fakepipelines",
		"cicd_pipelines",
		"cicd_scopes",
	}
	PluginConfigTables = []string{
		"_tool_fakeplugin_fakescopeconfigs",
		"_tool_fakeplugin_fakeconnections",
	}
	PluginScopeTable = "_tool_fakeplugin_fakeprojects"
)

type (
	FakePluginConnection struct {
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}
	FakeProject struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		ConnectionId  uint64 `json:"connectionId"`
		ScopeConfigId uint64 `json:"scopeConfigId"`
		Url           string `json:"url"`
	}
	FakeScopeConfig struct {
		Id       uint64   `json:"id"`
		Name     string   `json:"name"`
		Env      string   `json:"env"`
		Entities []string `json:"entities"`
	}
	BlueprintTestParams struct {
		connection *helper.Connection
		projects   []models.ApiOutputProject
		blueprints []models.Blueprint
		config     *FakeScopeConfig
		scope      *FakeProject
	}
)

func ConnectLocalServer(t *testing.T) *helper.DevlakeClient {
	fmt.Println("Connect to server")
	client := helper.StartDevLakeServer(t, nil)
	client.SetTimeout(0 * time.Second)
	client.SetPipelineTimeout(0 * time.Second)
	return client
}

func CreateClient(t *testing.T) *helper.DevlakeClient {
	path := filepath.Join(helper.ProjectRoot, FAKE_PLUGIN_DIR)
	_ = os.Setenv("REMOTE_PLUGIN_DIR", path)
	client := ConnectLocalServer(t)
	client.AwaitPluginAvailability(PLUGIN_NAME, 60*time.Second)
	return client
}

func CreateTestConnection(client *helper.DevlakeClient) *helper.Connection {
	connection := client.CreateConnection(PLUGIN_NAME,
		FakePluginConnection{
			Name:  "Test connection",
			Token: TOKEN,
		},
	)
	return connection
}

func CreateTestScope(client *helper.DevlakeClient, config *FakeScopeConfig, connectionId uint64) *FakeProject {
	scopes := helper.Cast[[]FakeProject](client.CreateScopes(PLUGIN_NAME,
		connectionId,
		FakeProject{
			Id:            "p1",
			Name:          "Project 1",
			ConnectionId:  connectionId,
			Url:           "http://fake.org/api/project/p1",
			ScopeConfigId: config.Id,
		},
	))
	return &scopes[0]
}

func CreateTestScopeConfig(client *helper.DevlakeClient, connectionId uint64) *FakeScopeConfig {
	config := helper.Cast[FakeScopeConfig](client.CreateScopeConfig(PLUGIN_NAME, connectionId, FakeScopeConfig{
		Name:     "Scope config",
		Env:      "test env",
		Entities: []string{plugin.DOMAIN_TYPE_CICD},
	}))
	return &config
}

func CreateTestBlueprints(t *testing.T, client *helper.DevlakeClient, count int) *BlueprintTestParams {
	t.Helper()
	connection := CreateTestConnection(client)
	config := CreateTestScopeConfig(client, connection.ID)
	scope := CreateTestScope(client, config, connection.ID)
	var bps []models.Blueprint
	var projects []models.ApiOutputProject
	for i := 1; i <= count; i++ {
		projectName := fmt.Sprintf("Test project %d", i)
		client.CreateProject(&helper.ProjectConfig{
			ProjectName: projectName,
		})
		project := client.GetProject(projectName)
		blueprint := client.PatchBasicBlueprintV2(
			project.Blueprint.ID,
			fmt.Sprintf("Test project %d-Blueprint", i),
			&helper.BlueprintV2Config{
				Connection: &models.BlueprintConnection{
					PluginName:   "fake",
					ConnectionId: connection.ID,
					Scopes: []*models.BlueprintScope{
						{
							ScopeId: scope.Id,
						},
					},
				},
				SkipOnFail:  true,
				ProjectName: projectName,
			},
		)
		bps = append(bps, blueprint)
		require.Equal(t, blueprint.Name, project.Blueprint.Name)
		projects = append(projects, project)
	}
	return &BlueprintTestParams{
		connection: connection,
		projects:   projects,
		blueprints: bps,
		config:     config,
		scope:      scope,
	}
}

func DeleteScopeWithDataIntegrityValidation(t *testing.T, client *helper.DevlakeClient, connectionId uint64, scopeId string, deleteDataOnly bool) services.BlueprintProjectPairs {
	db := client.GetDal()
	for _, table := range PluginDataTables {
		count, err := db.Count(dal.From(table))
		require.NoError(t, err)
		require.Greaterf(t, int(count), 0, fmt.Sprintf("no data was found in table: %s", table))
	}
	configData := map[string]int{}
	for _, table := range PluginConfigTables {
		count, err := db.Count(dal.From(table))
		require.NoError(t, err)
		require.Greaterf(t, int(count), 0, fmt.Sprintf("no data was found in table: %s", table))
		configData[table] = int(count)
	}
	refs := client.DeleteScope(PLUGIN_NAME, connectionId, scopeId, deleteDataOnly)
	if client.LastReturnedStatusCode() == http.StatusOK {
		for _, table := range PluginDataTables {
			count, err := db.Count(dal.From(table))
			require.NoError(t, err)
			require.Equalf(t, 0, int(count), fmt.Sprintf("data was found in table: %s", table))
		}
		count, err := db.Count(dal.From(PluginScopeTable))
		require.NoError(t, err)
		require.Equalf(t, 0, int(count), fmt.Sprintf("data was found in table: %s", PluginScopeTable))
	} else {
		for _, table := range PluginDataTables {
			count, err := db.Count(dal.From(table))
			require.NoError(t, err)
			require.Greaterf(t, int(count), 0, fmt.Sprintf("no data was found in table: %s", table))
		}
		count, err := db.Count(dal.From(PluginScopeTable))
		require.NoError(t, err)
		require.Greaterf(t, int(count), 0, fmt.Sprintf("no data was found in table: %s", PluginScopeTable))
	}
	for _, table := range PluginConfigTables {
		count, err := db.Count(dal.From(table))
		require.NoError(t, err)
		require.Equalf(t, configData[table], int(count), fmt.Sprintf("data was unexpectedly changed in table: %s", table))
	}
	return refs
}
