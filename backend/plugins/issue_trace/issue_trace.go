package main

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/plugins/issue_trace/impl"
	"github.com/spf13/cobra"
)

var PluginEntry impl.IssueTrace //nolint

func main() {
	cmd := &cobra.Command{Use: "spe"}
	connectionId := cmd.Flags().Uint64P("connection", "s", 0, "jira connection id")
	boardId := cmd.Flags().Uint64P("board", "b", 0, "jira board id")
	plugin := cmd.Flags().StringP("plugin", "p", "", "plugin")
	_ = cmd.MarkFlagRequired("connection")
	_ = cmd.MarkFlagRequired("board")
	_ = cmd.MarkFlagRequired("plugin")

	cmd.Run = func(c *cobra.Command, args []string) {
		fmt.Println("Direct running")
		runner.DirectRun(c, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"boardId":      *boardId,
			"plugin":       plugin,
		}, "")
	}
	runner.RunCmd(cmd)
}
