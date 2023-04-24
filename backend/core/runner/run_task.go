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
	gocontext "context"
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	contextimpl "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"time"
)

// RunTask FIXME ...
func RunTask(
	ctx gocontext.Context,
	basicRes context.BasicRes,
	progress chan plugin.RunningProgress,
	taskId uint64,
) (err errors.Error) {
	db := basicRes.GetDal()
	task := &models.Task{}
	if err := db.First(task, dal.Where("id = ?", taskId)); err != nil {
		return err
	}
	if task.Status == models.TASK_COMPLETED {
		return errors.Default.New("invalid task status")
	}
	dbPipeline := &models.Pipeline{}
	if err := db.First(dbPipeline, dal.Where("id = ? ", task.PipelineId)); err != nil {
		return err
	}
	logger, err := getTaskLogger(basicRes.GetLogger(), task)
	if err != nil {
		return err
	}
	beganAt := time.Now()
	// make sure task status always correct even if it panicked
	defer func() {
		if r := recover(); r != nil {
			var e error
			switch et := r.(type) {
			case error:
				e = et
			default:
				e = fmt.Errorf("%v", et)
			}
			err = errors.Default.Wrap(e, fmt.Sprintf("run task failed with panic (%s)", utils.GatherCallFrames(0)))
			logger.Error(err, "run task failed with panic")
		}
		finishedAt := time.Now()
		spentSeconds := finishedAt.Unix() - beganAt.Unix()
		if err != nil {
			lakeErr := errors.AsLakeErrorType(err)
			subTaskName := "unknown"
			if lakeErr = lakeErr.As(errors.SubtaskErr); lakeErr != nil {
				if meta, ok := lakeErr.GetData().(*plugin.SubTaskMeta); ok {
					subTaskName = meta.Name
				}
			} else {
				lakeErr = errors.Convert(err)
			}
			dbe := db.UpdateColumns(task, []dal.DalSet{
				{ColumnName: "status", Value: models.TASK_FAILED},
				{ColumnName: "message", Value: lakeErr.Error()},
				{ColumnName: "error_name", Value: lakeErr.Messages().Format()},
				{ColumnName: "finished_at", Value: finishedAt},
				{ColumnName: "spent_seconds", Value: spentSeconds},
				{ColumnName: "failed_sub_task", Value: subTaskName},
			})
			if dbe != nil {
				logger.Error(dbe, "failed to finalize task status into db (task failed)")
			}
		} else {
			dbe := db.UpdateColumns(task, []dal.DalSet{
				{ColumnName: "status", Value: models.TASK_COMPLETED},
				{ColumnName: "message", Value: ""},
				{ColumnName: "finished_at", Value: finishedAt},
				{ColumnName: "spent_seconds", Value: spentSeconds},
			})
			if dbe != nil {
				logger.Error(dbe, "failed to finalize task status into db (task succeeded)")
			}
		}
		// update finishedTasks
		dbe := db.UpdateColumn(
			&models.Pipeline{},
			"finished_tasks", dal.Expr("finished_tasks + 1"),
			dal.Where("id=?", task.PipelineId),
		)
		if dbe != nil {
			logger.Error(dbe, "update pipeline state failed")
		}
		// not return err if the `SkipOnFail` is true
		if dbPipeline.SkipOnFail {
			err = nil
		}
	}()

	// start execution
	logger.Info("start executing task: %d", task.ID)
	dbe := db.UpdateColumns(task, []dal.DalSet{
		{ColumnName: "status", Value: models.TASK_RUNNING},
		{ColumnName: "message", Value: ""},
		{ColumnName: "began_at", Value: beganAt},
	})
	if dbe != nil {
		return dbe
	}

	err = RunPluginTask(
		ctx,
		basicRes.ReplaceLogger(logger),
		task,
		progress,
	)
	return err
}

// RunPluginTask FIXME ...
func RunPluginTask(
	ctx gocontext.Context,
	basicRes context.BasicRes,
	task *models.Task,
	progress chan plugin.RunningProgress,
) errors.Error {
	pluginMeta, err := plugin.GetPlugin(task.Plugin)
	if err != nil {
		return errors.Default.WrapRaw(err)
	}
	pluginTask, ok := pluginMeta.(plugin.PluginTask)
	if !ok {
		return errors.Default.New(fmt.Sprintf("plugin %s doesn't support PluginTask interface", task.Plugin))
	}
	return RunPluginSubTasks(
		ctx,
		basicRes,
		task,
		pluginTask,
		progress,
	)
}

// RunPluginSubTasks FIXME ...
func RunPluginSubTasks(
	ctx gocontext.Context,
	basicRes context.BasicRes,
	task *models.Task,
	pluginTask plugin.PluginTask,
	progress chan plugin.RunningProgress,
) errors.Error {
	logger := basicRes.GetLogger()
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
	subtaskNames, err := task.GetSubTasks()
	if err != nil {
		return err
	}
	if len(subtaskNames) != 0 {
		// decode user specified subtasks
		var specifiedTasks []string
		err := api.Decode(subtaskNames, &specifiedTasks, nil)
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

	taskCtx := contextimpl.NewDefaultTaskContext(ctx, basicRes, task.Plugin, subtasksFlag, progress)
	if closeablePlugin, ok := pluginTask.(plugin.CloseablePluginTask); ok {
		defer closeablePlugin.Close(taskCtx)
	}
	options, err := task.GetOptions()
	if err != nil {
		return err
	}
	taskData, err := pluginTask.PrepareTaskData(taskCtx, options)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error preparing task data for %s", task.Plugin))
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
			progress <- plugin.RunningProgress{
				Type:          plugin.SetCurrentSubTask,
				SubTaskName:   subtaskMeta.Name,
				SubTaskNumber: subtaskNumber,
			}
		}
		err = runSubtask(basicRes, subtaskCtx, task.ID, subtaskNumber, subtaskMeta.EntryPoint)
		if err != nil {
			err = errors.SubtaskErr.Wrap(err, fmt.Sprintf("subtask %s ended unexpectedly", subtaskMeta.Name), errors.WithData(&subtaskMeta))
			logger.Error(err, "")
			return err
		}
		taskCtx.IncProgress(1)
	}

	return nil
}

// UpdateProgressDetail FIXME ...
func UpdateProgressDetail(basicRes context.BasicRes, taskId uint64, progressDetail *models.TaskProgressDetail, p *plugin.RunningProgress) {
	task := &models.Task{}
	task.ID = taskId
	switch p.Type {
	case plugin.TaskSetProgress:
		progressDetail.TotalSubTasks = p.Total
		progressDetail.FinishedSubTasks = p.Current
	case plugin.TaskIncProgress:
		progressDetail.FinishedSubTasks = p.Current
		// TODO: get rid of db update
		pct := float32(p.Current) / float32(p.Total)
		err := basicRes.GetDal().UpdateColumn(task, "progress", pct)
		if err != nil {
			basicRes.GetLogger().Error(err, "failed to update progress")
		}
	case plugin.SubTaskSetProgress:
		progressDetail.TotalRecords = p.Total
		progressDetail.FinishedRecords = p.Current
	case plugin.SubTaskIncProgress:
		progressDetail.FinishedRecords = p.Current
	case plugin.SetCurrentSubTask:
		progressDetail.SubTaskName = p.SubTaskName
		progressDetail.SubTaskNumber = p.SubTaskNumber
	}
}

func runSubtask(
	basicRes context.BasicRes,
	ctx plugin.SubTaskContext,
	parentID uint64,
	subtaskNumber int,
	entryPoint plugin.SubTaskEntryPoint,
) errors.Error {
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
		recordSubtask(basicRes, subtask)
	}()
	return entryPoint(ctx)
}

func recordSubtask(basicRes context.BasicRes, subtask *models.Subtask) {
	if err := basicRes.GetDal().Create(subtask); err != nil {
		basicRes.GetLogger().Error(err, "error writing subtask %d status to DB: %v", subtask.ID)
	}
}

func getTaskLogger(parentLogger log.Logger, task *models.Task) (log.Logger, errors.Error) {
	logger := parentLogger.Nested(fmt.Sprintf("task #%d", task.ID))
	loggingPath := logruslog.GetTaskLoggerPath(logger.GetConfig(), task)
	stream, err := logruslog.GetFileStream(loggingPath)
	if err != nil {
		return nil, err
	}
	logger.SetStream(&log.LoggerStreamConfig{
		Path:   loggingPath,
		Writer: stream,
	})
	return logger, nil
}
