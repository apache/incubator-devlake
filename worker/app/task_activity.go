package app

import (
	"context"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/runner"
	"go.temporal.io/sdk/activity"
)

func DevLakeTaskActivity(ctx context.Context, configJson []byte, taskId uint64) error {
	cfg, log, db, err := loadResources(configJson)
	log.Info("received task #%d", taskId)
	progress := make(chan core.RunningProgress)
	defer close(progress)
	go func() {
		for p := range progress {
			activity.RecordHeartbeat(ctx, p)
		}
	}()
	err = runner.RunTask(cfg, log, db, ctx, progress, taskId)
	if err != nil {
		log.Error("failed to execute task #%d: %w", taskId, err)
	}
	log.Info("finished task #%d", taskId)
	return err
}
