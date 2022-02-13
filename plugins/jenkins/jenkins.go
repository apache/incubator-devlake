package main

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	errors "github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jenkins/api"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"github.com/merico-dev/lake/plugins/jenkins/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

var _ core.Plugin = (*Jenkins)(nil)

type JenkinsOptions struct {
	Host     string
	Username string
	Password string
}

type Jenkins struct{}

func (j Jenkins) Init() {
	var err = lakeModels.Db.AutoMigrate(&models.JenkinsJob{}, &models.JenkinsBuild{})
	if err != nil {
		logger.Error("Failed to auto migrate jenkins models", err)
	}
}

func (j Jenkins) Description() string {
	return "Jenkins plugin"
}

func (j Jenkins) CleanData() {
	var err = lakeModels.Db.Exec("truncate table jenkins_jobs").Error
	if err != nil {
		logger.Error("Failed to truncate jenkins models", err)
	}
}

func (j Jenkins) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var op = JenkinsOptions{
		Host:     config.V.GetString("JENKINS_ENDPOINT"),
		Username: config.V.GetString("JENKINS_USERNAME"),
		Password: config.V.GetString("JENKINS_PASSWORD"),
	}

	var err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("Failed to decode options: %v", err)
	}

	var rateLimitPerSecondInt int
	rateLimitPerSecondInt, err = core.GetRateLimitPerSecond(options, 15)
	if err != nil {
		return err
	}

	scheduler, err := utils.NewWorkerScheduler(10, rateLimitPerSecondInt, ctx)
	defer scheduler.Release()
	if err != nil {
		return fmt.Errorf("could not create scheduler")
	}

	j.CleanData()
	var apiClient = tasks.CreateApiClient(nil, op.Host, op.Username, op.Password)

	progress <- float32(0.1)
	err = tasks.CollectJobs(apiClient, scheduler)
	if err != nil {
		logger.Error("Fail to sync jobs", err)
		return &errors.SubTaskError{
			SubTaskName: "CollectJobs",
			Message:     err.Error(),
		}
	}

	progress <- float32(0.4)
	err = tasks.CollectBuilds(apiClient, scheduler)
	if err != nil {
		logger.Error("Fail to collect builds", err)
		return &errors.SubTaskError{
			SubTaskName: "CollectJobs",
			Message:     err.Error(),
		}
	}

	progress <- float32(0.7)
	err = tasks.ConvertJobs()
	if err != nil {
		logger.Error("Fail to convert jobs", err)
		return &errors.SubTaskError{
			SubTaskName: "ConvertJobs",
			Message:     err.Error(),
		}
	}
	progress <- float32(0.8)
	err = tasks.ConvertBuilds()
	if err != nil {
		logger.Error("Fail to convert builds", err)
		return &errors.SubTaskError{
			SubTaskName: "ConvertBuilds",
			Message:     err.Error(),
		}
	}
	progress <- float32(1.0)

	return nil
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

	err := core.RegisterPlugin("jenkins", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
		close(progress)
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
