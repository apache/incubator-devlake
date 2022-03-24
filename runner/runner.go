package runner

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func RunPluginTask(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	options map[string]interface{},
	progress chan core.RunningProgress,
) error {
	pluginMeta, err := core.GetPlugin(name)
	if err != nil {
		return err
	}
	pluginTask, ok := pluginMeta.(core.PluginTask)
	if !ok {
		return fmt.Errorf("plugin %s doesn't support PluginTask interface", name)
	}

	return RunPluginSubTasks(
		cfg,
		logger,
		db,
		ctx,
		name,
		options,
		pluginTask,
		progress,
	)
}

func RunPluginSubTasks(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	options map[string]interface{},
	pluginTask core.PluginTask,
	progress chan core.RunningProgress,
) error {
	logger.Info("start plugin")

	// find out all possible subtasks this plugin can offer
	subtaskMetas := pluginTask.SubTaskMetas()
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
		if len(specifiedTasks) > 0 {
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
	}
	// make sure `Required` subtasks are always enabled
	for _, subtaskMeta := range subtaskMetas {
		if subtaskMeta.Required {
			subtasks[subtaskMeta.Name] = true
		}
	}
	// calculate total step(number of task to run)
	steps := 0
	for _, enabled := range subtasks {
		if enabled {
			steps++
		}
	}

	taskCtx := helper.NewDefaultTaskContext(cfg, logger, db, ctx, name, subtasks, progress)
	taskData, err := pluginTask.PrepareTaskData(taskCtx, options)
	if err != nil {
		return err
	}
	taskCtx.SetData(taskData)

	// execute subtasks in order
	taskCtx.SetProgress(0, steps)
	i := 0
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
		i++
		if progress != nil {
			progress <- core.RunningProgress{
				Type:          core.SetCurrentSubTask,
				SubTaskName:   subtaskMeta.Name,
				SubTaskNumber: i,
			}
		}
		err = subtaskMeta.EntryPoint(subtaskCtx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: subtaskMeta.Name,
				Message:     err.Error(),
			}
		}
		taskCtx.IncProgress(1)
	}

	return nil
}
