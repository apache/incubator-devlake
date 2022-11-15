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
	"sort"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// RunPipeline FIXME ...
func RunPipeline(
	_ *viper.Viper,
	log core.Logger,
	db *gorm.DB,
	pipelineId uint64,
	runTasks func([]uint64) errors.Error,
) errors.Error {
	// load tasks for pipeline
	var tasks []models.Task
	err := db.Where("pipeline_id = ? AND status = ?", pipelineId, models.TASK_CREATED).Order("pipeline_row, pipeline_col").Find(&tasks).Error
	if err != nil {
		return errors.Convert(err)
	}
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].PipelineRow < tasks[j].PipelineRow {
			return true
		} else if tasks[i].PipelineRow == tasks[j].PipelineRow {
			return tasks[i].PipelineCol < tasks[j].PipelineCol
		}
		return true
	})
	taskIds := make([][]uint64, 0)
	for _, task := range tasks {
		for len(taskIds) < task.PipelineRow {
			taskIds = append(taskIds, make([]uint64, 0))
		}
		taskIds[task.PipelineRow-1] = append(taskIds[task.PipelineRow-1], task.ID)
	}
	return runPipelineTasks(log, db, pipelineId, taskIds, runTasks)
}

func runPipelineTasks(
	log core.Logger,
	db *gorm.DB,
	pipelineId uint64,
	taskIds [][]uint64,
	runTasks func([]uint64) errors.Error,
) errors.Error {
	// load pipeline from db
	dbPipeline := &models.DbPipeline{}
	err := db.Find(dbPipeline, pipelineId).Error
	if err != nil {
		return errors.Convert(err)
	}

	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	for i, row := range taskIds {
		// update stage
		err = db.Model(dbPipeline).Updates(map[string]interface{}{
			"status": models.TASK_RUNNING,
			"stage":  i + 1,
		}).Error
		if err != nil {
			log.Error(err, "update pipeline state failed")
			break
		}
		// run tasks in parallel
		err = runTasks(row)
		if err != nil {
			log.Error(err, "run tasks failed")
			return errors.Convert(err)
		}

		// update finishedTasks
		err = db.Model(dbPipeline).Updates(map[string]interface{}{
			"finished_tasks": gorm.Expr("finished_tasks + ?", len(row)),
		}).Error
		if err != nil {
			log.Error(err, "update pipeline state failed")
			return errors.Convert(err)
		}
	}
	log.Info("pipeline finished in %d ms: %v", time.Now().UnixMilli()-dbPipeline.BeganAt.UnixMilli(), err)
	return errors.Convert(err)
}
