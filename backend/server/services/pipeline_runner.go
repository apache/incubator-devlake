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

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/worker/app"
	"go.temporal.io/sdk/client"
	"time"
)

type pipelineRunner struct {
	logger   log.Logger
	pipeline *models.Pipeline
}

func (p *pipelineRunner) runPipelineStandalone() errors.Error {
	return runner.RunPipeline(
		basicRes.ReplaceLogger(p.logger),
		p.pipeline.ID,
		func(taskIds []uint64) errors.Error {
			return RunTasksStandalone(p.logger, taskIds)
		},
	)
}

func (p *pipelineRunner) runPipelineViaTemporal() errors.Error {
	workflowOpts := client.StartWorkflowOptions{
		ID:        getTemporalWorkflowId(p.pipeline.ID),
		TaskQueue: cfg.GetString("TEMPORAL_TASK_QUEUE"),
	}
	// send only the very basis data
	configJson, err := json.Marshal(cfg.AllSettings())
	if err != nil {
		return errors.Convert(err)
	}
	p.logger.Info("enqueue pipeline #%d into temporal task queue", p.pipeline.ID)
	workflow, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOpts,
		app.DevLakePipelineWorkflow,
		configJson,
		p.pipeline.ID,
		p.logger.GetConfig(),
	)
	if err != nil {
		p.logger.Error(err, "failed to enqueue pipeline #%d into temporal", p.pipeline.ID)
		return errors.Convert(err)
	}
	err = workflow.Get(context.Background(), nil)
	if err != nil {
		p.logger.Info("failed to execute pipeline #%d via temporal: %v", p.pipeline.ID, err)
	}
	p.logger.Info("pipeline #%d finished by temporal", p.pipeline.ID)
	return errors.Convert(err)
}

// GetPipelineLogger returns logger for the pipeline
func GetPipelineLogger(pipeline *models.Pipeline) log.Logger {
	pipelineLogger := globalPipelineLog.Nested(
		fmt.Sprintf("pipeline #%d", pipeline.ID),
	)
	loggingPath := logruslog.GetPipelineLoggerPath(pipelineLogger.GetConfig(), pipeline)
	stream, err := logruslog.GetFileStream(loggingPath)
	if err != nil {
		globalPipelineLog.Error(nil, "unable to set stream for logging pipeline %d", pipeline.ID)
	} else {
		pipelineLogger.SetStream(&log.LoggerStreamConfig{
			Path:   loggingPath,
			Writer: stream,
		})
	}
	return pipelineLogger
}

// runPipeline start a pipeline actually
func runPipeline(pipelineId uint64) errors.Error {
	ppl, err := GetPipeline(pipelineId)
	if err != nil {
		return err
	}
	pipelineRun := pipelineRunner{
		logger:   GetPipelineLogger(ppl),
		pipeline: ppl,
	}
	// run
	if temporalClient != nil {
		err = pipelineRun.runPipelineViaTemporal()
	} else {
		err = pipelineRun.runPipelineStandalone()
	}
	isCancelled := errors.Is(err, context.Canceled)
	if err != nil {
		err = errors.Default.Wrap(err, fmt.Sprintf("Error running pipeline %d.", pipelineId))
	}
	dbPipeline, e := GetDbPipeline(pipelineId)
	if e != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("Unable to get pipeline %d.", pipelineId))
	}
	// finished, update database
	finishedAt := time.Now()
	dbPipeline.FinishedAt = &finishedAt
	if dbPipeline.BeganAt != nil {
		dbPipeline.SpentSeconds = int(finishedAt.Unix() - dbPipeline.BeganAt.Unix())
	}
	if err != nil {
		dbPipeline.Message = err.Error()
		dbPipeline.ErrorName = err.Messages().Format()
	}
	dbPipeline.Status, err = ComputePipelineStatus(dbPipeline, isCancelled)
	if err != nil {
		globalPipelineLog.Error(err, "compute pipeline status failed")
		return err
	}
	err = db.Update(dbPipeline)
	if err != nil {
		globalPipelineLog.Error(err, "update pipeline state failed")
		return err
	}
	// notify external webhook
	return NotifyExternal(pipelineId)
}

// ComputePipelineStatus determines pipleline status by its latest(rerun included) tasks statuses
// 1. TASK_COMPLETED: all tasks were executed sucessfully
// 2. TASK_FAILED: SkipOnFail=false with failed task(s)
// 3. TASK_PARTIAL: SkipOnFail=true with failed task(s)
func ComputePipelineStatus(pipeline *models.Pipeline, isCancelled bool) (string, errors.Error) {
	tasks, err := GetLatestTasksOfPipeline(pipeline)
	if err != nil {
		return "", err
	}

	succeeded, failed, pending, running := 0, 0, 0, 0

	for _, task := range tasks {
		if task.Status == models.TASK_COMPLETED {
			succeeded += 1
		} else if task.Status == models.TASK_FAILED || task.Status == models.TASK_CANCELLED {
			failed += 1
		} else if task.Status == models.TASK_RUNNING {
			running += 1
		} else {
			pending += 1
		}
	}

	if running > 0 || (!isCancelled && pipeline.SkipOnFail && pending > 0) {
		return "", errors.Default.New("unexpected status, did you call computePipelineStatus at a wrong timing?")
	}

	if failed == 0 {
		return models.TASK_COMPLETED, nil
	}
	if pipeline.SkipOnFail && succeeded > 0 {
		return models.TASK_PARTIAL, nil
	}
	return models.TASK_FAILED, nil
}

// GetLatestTasksOfPipeline returns latest tasks (reran tasks are excluding) of specified pipeline
func GetLatestTasksOfPipeline(pipeline *models.Pipeline) ([]*models.Task, errors.Error) {
	task := &models.Task{}
	cursor, err := db.Cursor(
		dal.From(task),
		dal.Where("pipeline_id = ?", pipeline.ID),
		dal.Orderby("id DESC"), // sort it by id so we can hit the latest task first for the RERUNed row/col
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	tasks := make([]*models.Task, 0, pipeline.TotalTasks)
	// define a struct for composite key to dedupe RERUNed tasks
	type rowcol struct{ row, col int }
	memorized := make(map[rowcol]bool)
	for cursor.Next() {
		if e := db.Fetch(cursor, task); e != nil {
			return nil, errors.Convert(e)
		}
		// dedupe reran tasks
		key := rowcol{task.PipelineRow, task.PipelineCol}
		if memorized[key] {
			continue
		}
		memorized[key] = true
		tasks = append(tasks, task)
		task = &models.Task{}
	}
	return tasks, nil
}
