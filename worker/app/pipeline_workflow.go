package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/merico-dev/lake/runner"
	"go.temporal.io/sdk/workflow"
)

func DevLakePipelineWorkflow(ctx workflow.Context, configJson []byte, pipelineId uint64) error {
	cfg, logger, db, err := loadResources(configJson)
	logger.Info("received pipeline #%d", pipelineId)
	err = runner.RunPipeline(
		cfg,
		logger,
		db,
		pipelineId,
		func(taskIds []uint64) error {
			futures := make([]workflow.Future, len(taskIds))
			for i, taskId := range taskIds {
				activityOpts := workflow.ActivityOptions{
					ActivityID:          fmt.Sprintf("task #%d", taskId),
					StartToCloseTimeout: 24 * time.Hour,
					WaitForCancellation: true,
				}
				activityCtx := workflow.WithActivityOptions(ctx, activityOpts)
				futures[i] = workflow.ExecuteActivity(activityCtx, DevLakeTaskActivity, configJson, taskId)
			}
			errs := make([]string, 0)
			for _, future := range futures {
				err := future.Get(ctx, nil)
				if err != nil {
					errs = append(errs, err.Error())
				}
			}
			if len(errs) > 0 {
				return fmt.Errorf(strings.Join(errs, "\n"))
			}
			return nil
		},
	)
	if err != nil {
		logger.Error("failed to execute pipeline #%d: %w", pipelineId, err)
	}
	logger.Info("finished pipeline #%d", pipelineId)
	return err
}
