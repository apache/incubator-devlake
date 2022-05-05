package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/plugins/gitlab/impl"
	"github.com/merico-dev/lake/runner"
	"github.com/spf13/cobra"
)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry impl.Gitlab //nolint

// standalone mode for debugging
func main() {
	gitlabCmd := &cobra.Command{Use: "gitlab"}
	projectId := gitlabCmd.Flags().IntP("project-id", "p", 0, "gitlab project id")

	_ = gitlabCmd.MarkFlagRequired("project-id")
	gitlabCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"projectId": *projectId,
		})
	}
	runner.RunCmd(gitlabCmd)
}
