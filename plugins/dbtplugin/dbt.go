package main

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/dbtplugin/tasks"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm" // A pseudo type for Plugin Interface implementation
)

var _ core.PluginMeta = (*Dbt)(nil)
var _ core.PluginInit = (*Dbt)(nil)
var _ core.PluginTask = (*Dbt)(nil)

type Dbt struct{}

func (plugin Dbt) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	// dbt: init
	return nil
}

func (plugin Dbt) Description() string {
	return "Convert data by dbt"
}

func (plugin Dbt) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.DbtConverterMeta,
	}
}

func (plugin Dbt) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.DbtOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	return &tasks.DbtTaskData{
		Options: &op,
	}, nil
}

func (plugin Dbt) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/dbtplugin"
}

var PluginEntry Dbt

// standalone mode for debugging
func main() {
	dbtCmd := &cobra.Command{Use: "dbt"}
	selectedModels := dbtCmd.Flags().StringP("models", "m", "my_first_dbt_model", "dbt select models")
	dbtCmd.MarkFlagRequired("selected_models")
	dbtCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"selectedModels": *selectedModels,
		})
	}
	runner.RunCmd(dbtCmd)
}
