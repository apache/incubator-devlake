package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/apache/incubator-devlake/worker/app"
	"go.temporal.io/sdk/client"
	"time"
)

type pipelineRunner struct {
	logger   core.Logger
	pipeline *models.Pipeline
}

func (p *pipelineRunner) runPipelineStandalone() error {
	return runner.RunPipeline(
		cfg,
		p.logger,
		db,
		p.pipeline.ID,
		func(taskIds []uint64) error {
			return runTasksStandalone(p.logger, taskIds)
		},
	)
}

func (p *pipelineRunner) runPipelineViaTemporal() error {
	workflowOpts := client.StartWorkflowOptions{
		ID:        getTemporalWorkflowId(p.pipeline.ID),
		TaskQueue: cfg.GetString("TEMPORAL_TASK_QUEUE"),
	}
	// send only the very basis data
	configJson, err := json.Marshal(cfg.AllSettings())
	if err != nil {
		return err
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
		p.logger.Error("failed to enqueue pipeline #%d into temporal", p.pipeline.ID)
		return err
	}
	err = workflow.Get(context.Background(), nil)
	if err != nil {
		p.logger.Info("failed to execute pipeline #%d via temporal: %w", p.pipeline.ID, err)
	}
	p.logger.Info("pipeline #%d finished by temporal", p.pipeline.ID)
	return err
}

func getPipelineLogger(pipeline *models.Pipeline) core.Logger {
	pipelineLogger := globalPipelineLog.Nested(
		fmt.Sprintf("pipeline #%d", pipeline.ID),
	)
	loggingPath := models.GetPipelineLoggerPath(pipelineLogger.GetConfig(), pipeline)
	pipelineLogger.GetConfig().Path = loggingPath
	if writer, err := pipelineLogger.GetConfig().GetStream(loggingPath); err != nil {
		globalPipelineLog.Error("unable to set stream for logging pipeline %d", pipeline.ID)
	} else {
		pipelineLogger.SetStream(writer)
	}
	return pipelineLogger
}

// runPipeline start a pipeline actually
func runPipeline(pipelineId uint64) error {
	pipeline, err := GetPipeline(pipelineId)
	if err != nil {
		return err
	}
	pipelineRun := pipelineRunner{
		logger:   getPipelineLogger(pipeline),
		pipeline: pipeline,
	}
	// run
	if temporalClient != nil {
		err = pipelineRun.runPipelineViaTemporal()
	} else {
		err = pipelineRun.runPipelineStandalone()
	}
	if err != nil {
		return fmt.Errorf("error running pipeline %d: %v", pipelineId, err)
	}
	pipeline, err = GetPipeline(pipelineId)
	if err != nil {
		return err
	}
	// finished, update database
	finishedAt := time.Now()
	pipeline.FinishedAt = &finishedAt
	pipeline.SpentSeconds = int(finishedAt.Unix() - pipeline.BeganAt.Unix())
	if err != nil {
		pipeline.Status = models.TASK_FAILED
		pipeline.Message = err.Error()
	} else {
		pipeline.Status = models.TASK_COMPLETED
		pipeline.Message = ""
	}
	dbe := db.Model(pipeline).Select("finished_at", "spent_seconds", "status", "message").Updates(pipeline).Error
	if dbe != nil {
		globalPipelineLog.Error("update pipeline state failed: %w", dbe)
		return dbe
	}
	// notify external webhook
	return NotifyExternal(pipelineId)
}
