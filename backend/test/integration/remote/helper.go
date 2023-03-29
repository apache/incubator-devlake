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
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/plugin"
	orgPlugin "github.com/apache/incubator-devlake/plugins/org/impl"
	"github.com/apache/incubator-devlake/test/integration/helper"
)

const (
	PLUGIN_NAME     = "fake"
	TOKEN           = "this_is_a_valid_token"
	FAKE_PLUGIN_DIR = "python/test/fakeplugin"
)

type (
	FakePluginConnection struct {
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}
	FakeProject struct {
		Id                   string `json:"id"`
		Name                 string `json:"name"`
		ConnectionId         uint64 `json:"connection_id"`
		TransformationRuleId uint64 `json:"transformation_rule_id"`
	}
	FakeTxRule struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
		Env  string `json:"env"`
	}
)

func SetupEnv() {
	fmt.Println("Setup test env")
	helper.Init()
	path := filepath.Join(helper.ProjectRoot, FAKE_PLUGIN_DIR, "start.sh")
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("REMOTE_PLUGINS_STARTUP_PATH", path)
	_ = os.Setenv("ENABLE_REMOTE_PLUGINS", "true")
}

func ConnectLocalServer(t *testing.T) *helper.DevlakeClient {
	fmt.Println("Connect to server")
	client := helper.ConnectLocalServer(t, &helper.LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		TruncateDb:   true,
		Plugins:      map[string]plugin.PluginMeta{"org": orgPlugin.Org{}},
	})
	client.SetTimeout(30 * time.Second)
	return client
}

func CreateClient(t *testing.T) *helper.DevlakeClient {
	SetupEnv()
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

func CreateTestScope(client *helper.DevlakeClient, connectionId uint64) any {
	res := client.CreateTransformationRule(PLUGIN_NAME, connectionId, FakeTxRule{Name: "Tx rule", Env: "test env"})
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
