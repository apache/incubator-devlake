package runner

import (
	"time"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RunPipeline(
	cfg *viper.Viper,
	log core.Logger,
	db *gorm.DB,
	pipelineId uint64,
	runTasks func([]uint64) error,
) error {
	// load pipeline from db
	pipeline := &models.Pipeline{}
	err := db.Find(pipeline, pipelineId).Error
	if err != nil {
		return err
	}
	// load tasks for pipeline
	var tasks []*models.Task
	err = db.Where("pipeline_id = ?", pipeline.ID).Order("pipeline_row, pipeline_col").Find(&tasks).Error
	if err != nil {
		return err
	}
	// convert to 2d array
	taskIds := make([][]uint64, 0)
	for _, task := range tasks {
		for len(taskIds) < task.PipelineRow {
			taskIds = append(taskIds, make([]uint64, 0))
		}
		taskIds[task.PipelineRow-1] = append(taskIds[task.PipelineRow-1], task.ID)
	}

	beganAt := time.Now()
	err = db.Model(pipeline).Updates(map[string]interface{}{
		"status":   models.TASK_RUNNING,
		"message":  "",
		"began_at": beganAt,
	}).Error
	if err != nil {
		return err
	}
	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	finishedTasks := 0
	for i, row := range taskIds {
		// update step
		err = db.Model(pipeline).Updates(map[string]interface{}{
			"status": models.TASK_RUNNING,
			"step":   i + 1,
		}).Error
		if err != nil {
			log.Error("update pipeline state failed: %w", err)
			break
		}
		// run tasks in parallel
		err = runTasks(row)
		if err != nil {
			log.Error("run tasks failed: %w", err)
			return err
		}
		// Deprecated
		// update finishedTasks
		finishedTasks += len(row)
		err = db.Model(pipeline).Updates(map[string]interface{}{
			"finished_tasks": finishedTasks,
		}).Error
		if err != nil {
			log.Error("update pipeline state failed: %w", err)
			return err
		}
	}

	log.Info("pipeline finished: %w", err)
	return err
}
