package main // must be main for plugin entry point

import (
	"fmt"

	"github.com/merico-dev/lake/migration"
	"github.com/merico-dev/lake/plugins/ae/api"
	"github.com/merico-dev/lake/plugins/ae/models/migrationscripts"
	"github.com/merico-dev/lake/plugins/ae/tasks"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*AE)(nil)
var _ core.PluginInit = (*AE)(nil)
var _ core.PluginTask = (*AE)(nil)
var _ core.PluginApi = (*AE)(nil)
var _ core.Migratable = (*AE)(nil)

type AE struct{}

func (plugin AE) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin AE) Description() string {
	return "To collect and enrich data from AE"
}

func (plugin AE) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectProjectMeta,
		tasks.CollectCommitsMeta,
		tasks.ExtractProjectMeta,
		tasks.ExtractCommitsMeta,
		tasks.ConvertCommitsMeta,
	}
}

func (plugin AE) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.AeOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.ProjectId <= 0 {
		return nil, fmt.Errorf("projectId is required")
	}
	apiClient, err := tasks.CreateApiClient(taskCtx)
	if err != nil {
		return nil, err
	}
	return &tasks.AeTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin AE) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/ae"
}

func (plugin AE) MigrationScripts() []migration.Script {
	return []migration.Script{new(migrationscripts.InitSchemas), new(migrationscripts.UpdateSchemas20220511)}
}

func (plugin AE) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
		"connections": {
			"GET": api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":   api.GetConnection,
			"PATCH": api.PatchConnection,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry AE //nolint

func main() {
	aeCmd := &cobra.Command{Use: "ae"}
	projectId := aeCmd.Flags().IntP("project-id", "p", 0, "ae project id")
	_ = aeCmd.MarkFlagRequired("project-id")
	aeCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"projectId": *projectId,
		})
	}
	runner.RunCmd(aeCmd)
}
