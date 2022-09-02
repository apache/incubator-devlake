/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// RunTask FIXME ...
func RunTask(
	ctx context.Context,
	_ *viper.Viper,
	parentLogger core.Logger,
	db *gorm.DB,
	progress chan core.RunningProgress,
	taskId uint64,
) error {
	task := &models.Task{}
	err := db.Find(task, taskId).Error
	if err != nil {
		return err
	}
	if task.Status == models.TASK_COMPLETED {
		return errors.Default.New("invalid task status")
	}
	log, err := getTaskLogger(parentLogger, task)
	if err != nil {
		return err
	}
	beganAt := time.Now()
	// make sure task status always correct even if it panicked
	defer func() {
		if r := recover(); r != nil {
			err = errors.Default.Wrap(r.(error), fmt.Sprintf("run task failed with panic (%s)", utils.GatherCallFrames(0)))
		}
		finishedAt := time.Now()
		spentSeconds := finishedAt.Unix() - beganAt.Unix()
		if err != nil {
			lakeErr := errors.AsLakeErrorType(err)
			subTaskName := "unknown"
			if lakeErr == nil {
				//skip
			} else if lakeErr = lakeErr.As(errors.SubtaskErr); lakeErr != nil {
				if meta, ok := lakeErr.GetData().(*core.SubTaskMeta); ok {
					subTaskName = meta.Name
				}
			}
			dbe := db.Model(task).Updates(map[string]interface{}{
				"status":          models.TASK_FAILED,
				"message":         err.Error(),
				"finished_at":     finishedAt,
				"spent_seconds":   spentSeconds,
				"failed_sub_task": subTaskName,
			}).Error
			if dbe != nil {
				log.Error(err, "failed to finalize task status into db")
			}
		} else {
			err = db.Model(task).Updates(map[string]interface{}{
				"status":        models.TASK_COMPLETED,
				"message":       "",
				"finished_at":   finishedAt,
				"spent_seconds": spentSeconds,
			}).Error
		}
	}()

	// start execution
	log.Info("start executing task: %d", task.ID)
	err = db.Model(task).Updates(map[string]interface{}{
		"status":   models.TASK_RUNNING,
		"message":  "",
		"began_at": beganAt,
	}).Error
	if err != nil {
		return err
	}

	var options map[string]interface{}
	err = json.Unmarshal(task.Options, &options)
	if err != nil {
		return err
	}
	var subtasks []string
	err = json.Unmarshal(task.Subtasks, &subtasks)
	if err != nil {
		return err
	}

	err = RunPluginTask(
		ctx,
		config.GetConfig(),
		log.Nested(task.Plugin),
		db,
		task.ID,
		task.Plugin,
		subtasks,
		options,
		progress,
	)
	return err
}

// RunPluginTask FIXME ...
func RunPluginTask(
	ctx context.Context,
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	taskID uint64,
	name string,
	subtasks []string,
	options map[string]interface{},
	progress chan core.RunningProgress,
) errors.Error {
	pluginMeta, err := core.GetPlugin(name)
	if err != nil {
		return errors.Default.WrapRaw(err)
	}
	pluginTask, ok := pluginMeta.(core.PluginTask)
	if !ok {
		return errors.Default.New(fmt.Sprintf("plugin %s doesn't support PluginTask interface", name))
	}
	return RunPluginSubTasks(
		ctx,
		cfg,
		logger,
		db,
		taskID,
		name,
		subtasks,
		options,
		pluginTask,
		progress,
	)
}

// RunPluginSubTasks FIXME ...
func RunPluginSubTasks(
	ctx context.Context,
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	taskID uint64,
	name string,
	subtaskNames []string,
	options map[string]interface{},
	pluginTask core.PluginTask,
	progress chan core.RunningProgress,
) errors.Error {
	logger.Info("start plugin")
	// find out all possible subtasks this plugin can offer
	subtaskMetas := pluginTask.SubTaskMetas()
	subtasksFlag := make(map[string]bool)
	for _, subtaskMeta := range subtaskMetas {
		subtasksFlag[subtaskMeta.Name] = subtaskMeta.EnabledByDefault
	}
	/* subtasksFlag example
	subtasksFlag := map[string]bool{
		"collectProject": true,
		"convertCommits": true,
		...
	}
	*/

	// user specifies what subtasks to run
	if len(subtaskNames) != 0 {
		// decode user specified subtasks
		var specifiedTasks []string
		err := mapstructure.Decode(subtaskNames, &specifiedTasks)
		if err != nil {
			return errors.Default.Wrap(err, "subtasks could not be decoded")
		}
		if len(specifiedTasks) > 0 {
			// first, disable all subtasks
			for task := range subtasksFlag {
				subtasksFlag[task] = false
			}
			// second, check specified subtasks is valid and enable them if so
			for _, task := range specifiedTasks {
				if _, ok := subtasksFlag[task]; ok {
					subtasksFlag[task] = true
				} else {
					return errors.Default.New(fmt.Sprintf("subtask %s does not exist", task))
				}
			}
		}
	}

	// make sure `Required` subtasks are always enabled
	for _, subtaskMeta := range subtaskMetas {
		if subtaskMeta.Required {
			subtasksFlag[subtaskMeta.Name] = true
		}
	}

	// calculate total step(number of task to run)
	steps := 0
	for _, enabled := range subtasksFlag {
		if enabled {
			steps++
		}
	}

	taskCtx := helper.NewDefaultTaskContext(ctx, cfg, logger, db, name, subtasksFlag, progress)
	if closeablePlugin, ok := pluginTask.(core.CloseablePluginTask); ok {
		defer closeablePlugin.Close(taskCtx)
	}
	taskData, err := pluginTask.PrepareTaskData(taskCtx, options)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error preparing task data for %s", name))
	}
	taskCtx.SetData(taskData)

	// execute subtasks in order
	taskCtx.SetProgress(0, steps)
	subtaskNumber := 0
	for _, subtaskMeta := range subtaskMetas {
		subtaskCtx, err := taskCtx.SubTaskContext(subtaskMeta.Name)
		if err != nil {
			// sth went wrong
			return errors.Default.Wrap(err, fmt.Sprintf("error getting context subtask %s", subtaskMeta.Name))
		}
		if subtaskCtx == nil {
			// subtask was disabled
			continue
		}

		// run subtask
		logger.Info("executing subtask %s", subtaskMeta.Name)
		subtaskNumber++
		if progress != nil {
			progress <- core.RunningProgress{
				Type:          core.SetCurrentSubTask,
				SubTaskName:   subtaskMeta.Name,
				SubTaskNumber: subtaskNumber,
			}
		}
		err = runSubtask(logger, db, taskID, subtaskNumber, subtaskCtx, subtaskMeta.EntryPoint)
		if err != nil {
			return errors.SubtaskErr.Wrap(err, fmt.Sprintf("subtask %s ended unexpectedly", subtaskMeta.Name), errors.WithData(&subtaskMeta))
		}
		taskCtx.IncProgress(1)
	}

	return nil
}

// UpdateProgressDetail FIXME ...
func UpdateProgressDetail(db *gorm.DB, logger core.Logger, taskId uint64, progressDetail *models.TaskProgressDetail, p *core.RunningProgress) {
	task := &models.Task{}
	task.ID = taskId
	switch p.Type {
	case core.TaskSetProgress:
		progressDetail.TotalSubTasks = p.Total
		progressDetail.FinishedSubTasks = p.Current
	case core.TaskIncProgress:
		progressDetail.FinishedSubTasks = p.Current
		// TODO: get rid of db update
		pct := float32(p.Current) / float32(p.Total)
		err := db.Model(task).Update("progress", pct).Error
		if err != nil {
			logger.Error(err, "failed to update progress: %w")
		}
	case core.SubTaskSetProgress:
		progressDetail.TotalRecords = p.Total
		progressDetail.FinishedRecords = p.Current
	case core.SubTaskIncProgress:
		progressDetail.FinishedRecords = p.Current
	case core.SetCurrentSubTask:
		progressDetail.SubTaskName = p.SubTaskName
		progressDetail.SubTaskNumber = p.SubTaskNumber
	}
}

func runSubtask(
	logger core.Logger,
	db *gorm.DB,
	parentID uint64,
	subtaskNumber int,
	ctx core.SubTaskContext,
	entryPoint core.SubTaskEntryPoint,
) error {
	beginAt := time.Now()
	subtask := &models.Subtask{
		Name:    ctx.GetName(),
		TaskID:  parentID,
		Number:  subtaskNumber,
		BeganAt: &beginAt,
	}
	defer func() {
		finishedAt := time.Now()
		subtask.FinishedAt = &finishedAt
		subtask.SpentSeconds = finishedAt.Unix() - beginAt.Unix()
		recordSubtask(logger, db, subtask)
	}()
	return entryPoint(ctx)
}

func recordSubtask(logger core.Logger, db *gorm.DB, subtask *models.Subtask) {
	if err := db.Create(&subtask).Error; err != nil {
		logger.Error(err, "error writing subtask %d status to DB: %v", subtask.ID)
	}
}

func getTaskLogger(parentLogger core.Logger, task *models.Task) (core.Logger, error) {
	log := parentLogger.Nested(fmt.Sprintf("task #%d", task.ID))
	loggingPath := logger.GetTaskLoggerPath(log.GetConfig(), task)
	stream, err := logger.GetFileStream(loggingPath)
	if err != nil {
		return nil, err
	}
	log.SetStream(&core.LoggerStreamConfig{
		Path:   loggingPath,
		Writer: stream,
	})
	return log, nil
}
