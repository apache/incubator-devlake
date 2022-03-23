package main

import (
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jenkins/api"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"github.com/merico-dev/lake/plugins/jenkins/tasks"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Jenkins)(nil)
var _ core.PluginInit = (*Jenkins)(nil)
var _ core.PluginTask = (*Jenkins)(nil)
var _ core.PluginApi = (*Jenkins)(nil)

type Jenkins struct{}

func (plugin Jenkins) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return db.AutoMigrate(
		&models.JenkinsJob{},
		&models.JenkinsBuild{},
	)
}

func (plugin Jenkins) Description() string {
	return "To collect and enrich data from Jenkins"
}

func (plugin Jenkins) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.CollectApiBuildsMeta,
		tasks.ExtractApiBuildsMeta,
		tasks.ConvertJobsMeta,
		tasks.ConvertBuildsMeta,
	}
}
func (plugin Jenkins) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.JenkinsOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode options: %v", err)
	}
	apiClient, err := tasks.CreateApiClient(taskCtx)
	if err != nil {
		return nil, err
	}
	return &tasks.JenkinsTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jenkins"
}

func (plugin Jenkins) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"sources": {
			"GET":  api.ListSources,
			"POST": api.PostSource,
		},
		"sources/:sourceId": {
			"GET": api.GetSource,
			"PUT": api.PutSource,
		},
	}
}

var PluginEntry Jenkins //nolint

func main() {
	jenkinsCmd := &cobra.Command{Use: "jenkins"}
	jenkinsCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{})
	}
	runner.RunCmd(jenkinsCmd)
}
