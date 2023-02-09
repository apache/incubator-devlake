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
	"testing"

	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/server/services/remote"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"github.com/apache/incubator-devlake/server/services/remote/models"
	plg "github.com/apache/incubator-devlake/server/services/remote/plugin"
)

type CircleCIConnection struct {
	ID    uint64 `json:"id"`
	Token string `json:"token" encrypt:"yes"`
}

func CreateRemotePlugin(t *testing.T) models.RemotePlugin {
	// TODO: Create a dummy plugin for tests instead of using CircleCI plugin
	pluginCmdPath := "../../plugins/circle_ci/circle_ci/main.py"
	invoker := bridge.NewPythonPoetryCmdInvoker(pluginCmdPath)

	pluginInfo := models.PluginInfo{}
	err := invoker.Call("plugin-info", bridge.DefaultContext).Get(&pluginInfo)

	if err != nil {
		t.Error("Cannot get plugin info", err)
		return nil
	}

	remotePlugin, err := remote.NewRemotePlugin(&pluginInfo)
	if err != nil {
		t.Error("Cannot create remote plugin", err)
		return nil
	}

	return remotePlugin
}

func TestCreateRemotePlugin(t *testing.T) {
	_ = CreateRemotePlugin(t)
}

func TestRunSubTask(t *testing.T) {
	remotePlugin := CreateRemotePlugin(t)
	dataflowTester := e2ehelper.NewDataFlowTester(t, "circleci", remotePlugin)
	subtask := remotePlugin.SubTaskMetas()[0]
	options := make(map[string]interface{})
	options["project_slug"] = "gh/circleci/bond"
	options["scopeId"] = "1"
	taskData := plg.RemotePluginTaskData{
		DbUrl:      bridge.DefaultContext.GetConfig("db_url"),
		Connection: CircleCIConnection{ID: 1},
		Options:    options,
	}
	dataflowTester.Subtask(subtask, taskData)
}

func TestTestConnection(t *testing.T) {
	remotePlugin := CreateRemotePlugin(t)

	var handler plugin.ApiResourceHandler
	for resource, endpoints := range remotePlugin.ApiResources() {
		if resource == "test" {
			handler = endpoints["POST"]
		}
	}

	if handler == nil {
		t.Error("Missing test connection API resource")
	}

	input := plugin.ApiResourceInput{
		Body: map[string]interface{}{
			"token": "secret",
		},
	}
	_, err := handler(&input)
	if err != nil {
		t.Error(err)
	}
}
