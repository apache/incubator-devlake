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

package runner

import (
	"context"
	"os"
	"os/exec"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/cobra"
)

func RunCmd(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("tasks", "t", nil, "specify what tasks to run, --tasks=collectIssues,extractIssues")
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

// DirectRun direct run plugin from command line.
// cmd: type is cobra.Command
// args: command line arguments
// pluginTask: specific built-in plugin, for example: feishu, jira...
// options: plugin config
func DirectRun(cmd *cobra.Command, args []string, pluginTask core.PluginTask, options map[string]interface{}) {
	tasks, err := cmd.Flags().GetStringSlice("tasks")
	if err != nil {
		panic(err)
	}
	options["tasks"] = tasks
	cfg := config.GetConfig()
	log := logger.Global.Nested(cmd.Use)
	db, err := NewGormDb(cfg, log)
	if err != nil {
		panic(err)
	}
	if pluginInit, ok := pluginTask.(core.PluginInit); ok {
		err = pluginInit.Init(cfg, log, db)
		if err != nil {
			panic(err)
		}
	}
	err = core.RegisterPlugin(cmd.Use, pluginTask.(core.PluginMeta))
	if err != nil {
		panic(err)
	}

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		} else if n == 1 && buf[0] == 99 {
			cancel()
		} else {
			println("unknown key press, code: ", buf[0])
		}
	}()
	println("press `c` to send cancel signal")

	err = RunPluginSubTasks(
		cfg,
		log,
		db,
		ctx,
		cmd.Use,
		options,
		pluginTask,
		nil,
	)
	if err != nil {
		panic(err)
	}
}
