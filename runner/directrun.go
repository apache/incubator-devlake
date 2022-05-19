package runner

import (
	"context"
	"os"
	"os/exec"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
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
