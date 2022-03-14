package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/ae/api"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/ae/tasks"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

var _ core.Plugin = (*AE)(nil)
var _ core.ManagedSubTasks = (*AE)(nil)

type AE string

func (plugin AE) Init() {
	logger.Info("INFO >>> init go plugin", true)
	err := lakeModels.Db.AutoMigrate(
		&models.AEProject{},
		&models.AECommit{})
	if err != nil {
		logger.Error("Error migrating ae: ", err)
		panic(err)
	}
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
	return &tasks.AeTaskData{
		Options:   &op,
		ApiClient: tasks.CreateApiClient(taskCtx.GetContext()),
	}, nil
}

func (plugin AE) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	logger := helper.NewDefaultTaskLogger(nil, "ae")
	logger.Info("start plugin")

	// find out all possible subtasks this plugin can offer
	subtaskMetas := plugin.SubTaskMetas()
	subtasks := make(map[string]bool)
	for _, subtaskMeta := range subtaskMetas {
		subtasks[subtaskMeta.Name] = subtaskMeta.EnabledByDefault
	}
	/* subtasks example
	subtasks := map[string]bool{
		"collectProject": true,
		"convertCommits": true,
		...
	}
	*/

	// if user specified what subtasks to run, obey
	// TODO: move tasks field to outer level, and rename it to subtasks
	if _, ok := options["tasks"]; ok {
		// decode user specified subtasks
		var specifiedTasks []string
		err := mapstructure.Decode(options["tasks"], &specifiedTasks)
		if err != nil {
			return err
		}
		// first, disable all subtasks
		for task := range subtasks {
			subtasks[task] = false
		}
		// second, check specified subtasks is valid and enable them if so
		for _, task := range specifiedTasks {
			if _, ok := subtasks[task]; ok {
				subtasks[task] = true
			} else {
				return fmt.Errorf("subtask %s does not exist", task)
			}
		}
	}

	taskCtx := helper.NewDefaultTaskContext("ae", ctx, logger, nil, subtasks)
	taskData, err := plugin.PrepareTaskData(taskCtx, options)
	if err != nil {
		return err
	}
	taskCtx.SetData(taskData)

	// execute subtasks in order
	for _, subtaskMeta := range subtaskMetas {
		subtaskCtx, err := taskCtx.SubTaskContext(subtaskMeta.Name)
		if err != nil {
			// sth went wrong
			return err
		}
		if subtaskCtx == nil {
			// subtask was disabled
			continue
		}

		// run subtask
		logger.Info("executing subtask %s", subtaskMeta.Name)
		err = subtaskMeta.EntryPoint(subtaskCtx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: subtaskMeta.Name,
				Message:     err.Error(),
			}
		}
	}

	return nil
}

func (plugin AE) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/ae"
}

func (plugin AE) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
		"sources": {
			"GET":  api.ListSources,
			"POST": api.PutSource,
		},
		"sources/:sourceId": {
			"GET": api.GetSource,
			"PUT": api.PutSource,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry AE //nolint

// standalone mode for debugging
func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic(fmt.Errorf("Usage: ae <project_id> [sub_task_name]"))
	}
	projectId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		panic(fmt.Errorf("error paring project_id: %w", err))
	}
	options := map[string]interface{}{
		"projectId": projectId,
	}
	if len(args) > 1 {
		options["tasks"] = []string{args[1]}
	}

	err = core.RegisterPlugin("ae", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			options,
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
