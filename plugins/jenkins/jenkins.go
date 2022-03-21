package main

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/helper"
	"time"

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

type Jenkins struct{}

func (plugin Jenkins) Init() {
	var err = lakeModels.Db.AutoMigrate(
		&models.JenkinsJob{},
		&models.JenkinsBuild{})
	if err != nil {
		logger.Error("Failed to auto migrate jenkins models", err)
	}
}

func (plugin Jenkins) Description() string {
	return "Jenkins plugin"
}

func (plugin Jenkins) CleanData() {
	var err = lakeModels.Db.Exec("truncate table jenkins_jobs").Error
	if err != nil {
		logger.Error("Failed to truncate jenkins models", err)
	}
}

func (plugin Jenkins) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var op tasks.JenkinsOptions
	var err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("Failed to decode options: %v", err)
	}
	v := config.GetConfig()

	var since *time.Time
	if op.Since != "" {
		*since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return fmt.Errorf("invalid value for `since`: %w", err)
		}
	}
	logger := helper.NewDefaultTaskLogger(nil, "jenkins")

	tasksToRun := make(map[string]bool, len(op.Tasks))
	if len(op.Tasks) == 0 {
		tasksToRun = map[string]bool{
			"collectApiJobs":   false,
			"extractApiJobs":   false,
			"collectApiBuilds": false,
			"extractApiBuilds": false,
			"convertJobs":      false,
			"convertBuilds":    false,
		}
	} else {
		for _, task := range op.Tasks {
			tasksToRun[task] = true
		}
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

	plugin.CleanData()

	op = tasks.JenkinsOptions{
		Host:     v.GetString("JENKINS_ENDPOINT"),
		Username: v.GetString("JENKINS_USERNAME"),
		Password: v.GetString("JENKINS_PASSWORD"),
	}

	var apiClient = tasks.NewJenkinsApiClient(op.Host, op.Username, op.Password, "", ctx, scheduler, logger)

	taskData := &tasks.JenkinsTaskData{
		Options:   &op,
		ApiClient: &apiClient.ApiClient,
		Since:     since,
	}

	taskCtx := helper.NewDefaultTaskContext("github", ctx, logger, taskData, tasksToRun)
	newTasks := []struct {
		name       string
		entryPoint core.SubTaskEntryPoint
	}{
		//{name: "collectApiJobs", entryPoint: tasks.CollectApiJobs},
		//{name: "extractApiJobs", entryPoint: tasks.ExtractApiJobs},
		{name: "collectApiBuilds", entryPoint: tasks.CollectApiBuilds},
		//{name: "extractApiBuilds", entryPoint: tasks.ExtractApiBuilds},
		//{name: "convertJobs", entryPoint: tasks.ConvertJobs},
		//{name: "convertBuilds", entryPoint: tasks.ConvertBuilds},
	}
	for _, t := range newTasks {
		c, err := taskCtx.SubTaskContext(t.name)
		if err != nil {
			return err
		}
		if c != nil {
			err = t.entryPoint(c)
			if err != nil {
				return &errors.SubTaskError{
					SubTaskName: t.name,
					Message:     err.Error(),
				}
			}
		}
	}

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
			map[string]interface{}{
				"tasks": []string{
					"collectApiJobs",
					"extractApiJobs",
					"collectApiBuilds",
					"extractApiBuilds",
					"convertJobs",
					"convertBuilds",
				},
			},
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
