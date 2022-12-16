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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/apache/incubator-devlake/worker/app"
	"go.temporal.io/sdk/client"
)

type pipelineRunner struct {
	logger   core.Logger
	pipeline *models.Pipeline
}

func (p *pipelineRunner) runPipelineStandalone() errors.Error {
	return runner.RunPipeline(
		runner.CreateBasicRes(cfg, p.logger, db),
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
func GetPipelineLogger(pipeline *models.Pipeline) core.Logger {
	pipelineLogger := globalPipelineLog.Nested(
		fmt.Sprintf("pipeline #%d", pipeline.ID),
	)
	loggingPath := logger.GetPipelineLoggerPath(pipelineLogger.GetConfig(), pipeline)
	stream, err := logger.GetFileStream(loggingPath)
	if err != nil {
		globalPipelineLog.Error(nil, "unable to set stream for logging pipeline %d", pipeline.ID)
	} else {
		pipelineLogger.SetStream(&core.LoggerStreamConfig{
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
		dbPipeline.Status = models.TASK_FAILED
		if lakeErr := errors.AsLakeErrorType(err); lakeErr != nil {
			dbPipeline.Message = lakeErr.Messages().Format()
		} else {
			dbPipeline.Message = err.Error()
		}
	} else {
		dbPipeline.Status = models.TASK_COMPLETED
		dbPipeline.Message = ""
	}
	dbe := db.Model(dbPipeline).Updates(dbPipeline).Error
	if dbe != nil {
		globalPipelineLog.Error(dbe, "update pipeline state failed")
		return errors.Convert(dbe)
	}
	// notify external webhook
	return NotifyExternal(pipelineId)
}
