package main

import (
	"github.com/merico-dev/lake/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/feishu/models"
	"github.com/merico-dev/lake/plugins/feishu/tasks"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm" // A pseudo type for Plugin Interface implementation
)

var _ core.PluginMeta = (*Feishu)(nil)
var _ core.PluginInit = (*Feishu)(nil)
var _ core.PluginTask = (*Feishu)(nil)
var _ core.PluginApi = (*Feishu)(nil)

type Feishu struct{}

func (plugin Feishu) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error{
	// feishu: init
	return db.AutoMigrate(
		&models.FeishuMeetingTopUserItem{},
	)
}

func (plugin Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (plugin Feishu) SubTaskMetas() []core.SubTaskMeta{
	return []core.SubTaskMeta{
		tasks.CollectMeetingTopUserItemMeta,
		tasks.ExtractMeetingTopUserItemMeta,
	}
}

func (plugin Feishu) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error){
	var op tasks.FeishuOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	apiClient, err := tasks.NewFeishuApiClient(taskCtx)
	if err != nil{
		return nil, err
	}
	return &tasks.FeishuTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Feishu) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/feishu"
}

func (plugin Feishu) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{}
}

var PluginEntry Feishu

// standalone mode for debugging
func main(){
	feishuCmd := &cobra.Command{Use: "feishu"}
	numOfDaysToCollect := feishuCmd.Flags().IntP("numOfDaysToCollect", "n", 8, "feishu collect days")
	feishuCmd.MarkFlagRequired("numOfDaysToCollect")
	feishuCmd.Run = func(cmd *cobra.Command, args []string){
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"numOfDaysToCollect": *numOfDaysToCollect,
		})
	}
	runner.RunCmd(feishuCmd)
}
