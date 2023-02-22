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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/require"
)

const PLUGIN_NAME = "fake"

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

func TestRunPipeline(t *testing.T) {
	setupEnv()
	buildPython(t)
	client := connectLocalServer(t)
	fmt.Println("Create new connection")
	conn := client.CreateConnection(PLUGIN_NAME,
		api.AccessToken{
			Token: "this_is_a_valid_token",
		},
	)
	client.SetTimeout(0)
	conns := client.ListConnections(PLUGIN_NAME)
	require.Equal(t, 1, len(conns))
	require.Equal(t, "this_is_a_valid_token", conns[0].Token)
	fmt.Println("Run pipeline")
	t.Run("run_pipeline", func(t *testing.T) {
		pipeline := client.RunPipeline(models.NewPipeline{
			Name: "remote_test",
			Plan: []plugin.PipelineStage{
				{
					{
						Plugin:   PLUGIN_NAME,
						Subtasks: nil,
						Options: map[string]interface{}{
							"connectionId": conn.ID,
							"scopeId":      "org/project",
						},
					},
				},
			},
		})
		require.Equal(t, models.TASK_COMPLETED, pipeline.Status)
		require.Equal(t, 1, pipeline.FinishedTasks)
		require.Equal(t, "", pipeline.ErrorName)
	})
}
