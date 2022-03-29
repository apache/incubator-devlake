package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RunPipeline(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	pipelineId uint64,
	runTask func(uint64) error,
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
	rowResults := make(chan error)
	rowErrors := make([]string, 0)
	for _, row := range taskIds {
		rowFinished := 0
		for _, taskId := range row {
			taskId := taskId
			go func() {
				logger.Info("run task in background ", taskId)
				rowResults <- runTask(taskId)
			}()
		}
		for err = range rowResults {
			finishedTasks++
			rowFinished++
			if err != nil {
				logger.Error("pipeline task failed", err)
				rowErrors = append(rowErrors, err.Error())
			}
			err = db.Model(pipeline).Updates(map[string]interface{}{
				"status":         models.TASK_RUNNING,
				"finished_tasks": finishedTasks,
			}).Error
			if err != nil {
				logger.Error("update pipeline state failed", err)
				rowErrors = append(rowErrors, err.Error())
			}
			if rowFinished == len(row) {
				break
			}
		}
		if len(rowErrors) > 0 {
			err = fmt.Errorf(strings.Join(rowErrors, "\n"))
			break
		}
	}
	close(rowResults)

	logger.Info("pipeline finished:", err == nil)
	// finished, update database
	finishedAt := time.Now()
	spentSeconds := finishedAt.Unix() - beganAt.Unix()
	if err != nil {
		err = db.Model(pipeline).Updates(map[string]interface{}{
			"status":        models.TASK_FAILED,
			"message":       err.Error(),
			"finished_at":   finishedAt,
			"spent_seconds": spentSeconds,
		}).Error
	} else {
		err = db.Model(pipeline).Updates(map[string]interface{}{
			"status":        models.TASK_COMPLETED,
			"message":       "",
			"finished_at":   finishedAt,
			"spent_seconds": spentSeconds,
		}).Error
	}
	return err
}
