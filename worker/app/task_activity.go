package app

import (
	"context"

	"github.com/merico-dev/lake/runner"
)

func DevLakeTaskActivity(ctx context.Context, configJson []byte, taskId uint64) error {
	cfg, logger, db, err := loadResources(configJson)
	logger.Info("received task #%d", taskId)
	err = runner.RunTask(cfg, logger, db, ctx, nil, taskId)
	if err != nil {
		logger.Error("failed to execute task #%d: %w", taskId, err)
	}
	logger.Info("finished task #%d", taskId)
	return err
}
