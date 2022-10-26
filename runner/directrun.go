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
	goerror "errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/impl"
	"github.com/apache/incubator-devlake/impl/dalgorm"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/cobra"
)

// RunCmd FIXME ...
func RunCmd(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("subtasks", "t", nil, "specify what tasks to run, --subtasks=collectIssues,extractIssues")
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
	tasks, err := cmd.Flags().GetStringSlice("subtasks")
	if err != nil {
		panic(err)
	}
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

	// collect migration and run
	migrator, err := InitMigrator(impl.NewDefaultBasicRes(cfg, log, dalgorm.NewDalgorm(db)))
	if err != nil {
		panic(err)
	}
	if migratable, ok := pluginTask.(core.PluginMigration); ok {
		migrator.Register(migratable.MigrationScripts(), cmd.Use)
	}
	err = migrator.Execute()
	if err != nil {
		panic(err)
	}
	ctx := createContext()
	var parentTaskID uint64
	err = RunPluginSubTasks(
		ctx,
		cfg,
		log,
		db,
		parentTaskID,
		cmd.Use,
		tasks,
		options,
		pluginTask,
		nil,
	)
	if err != nil {
		panic(err)
	}
}

func createContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, getStopSignals()...)
	go func() {
		<-sigc
		cancel()
	}()
	go func() {
		var buf string

		n, err := fmt.Scan(&buf)
		if err != nil {
			if goerror.Is(err, io.EOF) {
				return
			}
			panic(err)
		} else if n == 1 && buf == "c" {
			cancel()
		} else {
			println("unknown key press, code: ", buf)
			println("press `c` and enter to send cancel signal")
		}
	}()
	println("press `c` and enter to send cancel signal")
	return ctx
}

func getStopSignals() []os.Signal {
	if runtime.GOOS == "windows" {
		return []os.Signal{
			syscall.Signal(0x6), //syscall.SIGABRT for windows
		}
	}
	return []os.Signal{
		syscall.Signal(0x14), //syscall.SIGTSTP for posix
	}
}
