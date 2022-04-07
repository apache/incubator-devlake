package main

import (
	"github.com/merico-dev/lake/plugins/refdiff/tasks"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/merico-dev/lake/plugins/core"
)

// make sure interface is implemented
var _ core.PluginMeta = (*RefDiff)(nil)
var _ core.PluginInit = (*RefDiff)(nil)
var _ core.PluginTask = (*RefDiff)(nil)
var _ core.PluginApi = (*RefDiff)(nil)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry RefDiff //nolint

type RefDiff struct{}

func (plugin RefDiff) Description() string {
	return "Calculate commits diff for specified ref pairs based on `commits` and `commit_parents` tables"
}

func (plugin RefDiff) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin RefDiff) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		//tasks.CalculateCommitsDiffMeta,
		tasks.CalculateIssuesDiffMeta,
		tasks.CalculatePrCherryPickMeta,
	}
}

func (plugin RefDiff) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.RefdiffOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	return &tasks.RefdiffTaskData{
		Options: &op,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin RefDiff) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/refdiff"
}

func (plugin RefDiff) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}

// standalone mode for debugging
func main() {
	refdiffCmd := &cobra.Command{Use: "refdiff"}

	repoId := refdiffCmd.Flags().StringP("repo-id", "r", "", "repo id")
	newRef := refdiffCmd.Flags().StringP("new-ref", "n", "", "new ref")
	oldRef := refdiffCmd.Flags().StringP("old-ref", "o", "", "old ref")

	_ = refdiffCmd.MarkFlagRequired("repo-id")
	_ = refdiffCmd.MarkFlagRequired("new-ref")
	_ = refdiffCmd.MarkFlagRequired("old-ref")

	pairs := []map[string]string{
		{
			"NewRef": *newRef,
			"OldRef": *oldRef,
		},
	}
	refdiffCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"repoId": repoId,
			"pairs":  pairs,
		})
	}
	runner.RunCmd(refdiffCmd)
}
