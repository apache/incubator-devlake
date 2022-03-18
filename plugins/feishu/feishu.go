package main

import (
	"fmt"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/helper"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/feishu/models"
	"github.com/merico-dev/lake/plugins/feishu/tasks"
	"github.com/merico-dev/lake/utils"
)

var _ core.Plugin = (*Feishu)(nil)

type Feishu string


func (plugin Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (plugin Feishu) Init() {
	err := lakeModels.Db.AutoMigrate(
		&models.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin Feishu) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {

	var op tasks.FeishuOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil{
		return err
	}

	rateLimitPerSecondInt, err := core.GetRateLimitPerSecond(options, 5)
	if err != nil {
		return err
	}
	
	scheduler, err := utils.NewWorkerScheduler(10, rateLimitPerSecondInt, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	var FEISHU_ENDPOINT = config.GetConfig().GetString("FEISHU_ENDPOINT")
	// prepare contextual variables
	logger := helper.NewDefaultTaskLogger(nil, "feishu")
	apiClient := tasks.NewFeishuApiClient(
		FEISHU_ENDPOINT,
		scheduler,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create feishu api client: %w", err)
	}
	taskData := &tasks.FeishuTaskData{
		Options: &op,
		ApiClient: &apiClient.ApiClient,
	}

	tasksToRun := make(map[string]bool, len(op.Tasks))
	
	if len(op.Tasks) == 0{
		tasksToRun = map[string]bool{
			"collectMeetingTopUserItem": true,
			"extractMeetingTopUserItem": true,
		}
	}else{
		for _, task := range op.Tasks{
			tasksToRun[task] = true
		}
	}
	taskCtx := helper.NewDefaultTaskContext("feishu", ctx, logger, taskData, tasksToRun)
	newTasks := []struct{
		name string
		entryPoint core.SubTaskEntryPoint
	}{
		{name: "collectMeetingTopUserItem", entryPoint: tasks.CollectMeetingTopUserItem},
		{name: "extractMeetingTopUserItem", entryPoint: tasks.ExtractMeetingTopUserItem},
	}
	
	for _, t := range newTasks{
		c, err := taskCtx.SubTaskContext(t.name)
		if err != nil{
			return err
		}
		if c != nil {
			err = t.entryPoint(c)
			if err != nil{
				return &errors.SubTaskError{
					SubTaskName: t.name,
					Message: err.Error(),
				}
			}
		}
	}
	logger.Info("feishu plugin is end")
	return nil

}

func (plugin Feishu) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/feishu"
}

func (plugin Feishu) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{}
}

var PluginEntry Feishu

// standalone mode for debugging
func main() {
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"numOfDaysToCollect": 80,
			},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
