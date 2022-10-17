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
	"fmt"
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
	startTime := time.Now()
	// load pipeline from db
	dbPipeline := &models.DbPipeline{}
	err := db.Find(dbPipeline, pipelineId).Error
	if err != nil {
		return errors.Convert(err)
	}
	// load tasks for pipeline
	var tasks []*models.Task
	err = db.Where("pipeline_id = ?", dbPipeline.ID).Order("pipeline_row, pipeline_col").Find(&tasks).Error
	if err != nil {
		return errors.Convert(err)
	}
	if len(tasks) != dbPipeline.TotalTasks {
		return errors.Internal.New(fmt.Sprintf("expected total tasks to be %v, got %v", dbPipeline.TotalTasks, len(tasks)))
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
	err = db.Model(dbPipeline).Updates(map[string]interface{}{
		"status":   models.TASK_RUNNING,
		"message":  "",
		"began_at": beganAt,
	}).Error
	if err != nil {
		return errors.Convert(err)
	}
	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	finishedTasks := 0
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
		// Deprecated
		// update finishedTasks
		finishedTasks += len(row)
		err = db.Model(dbPipeline).Updates(map[string]interface{}{
			"finished_tasks": finishedTasks,
		}).Error
		if err != nil {
			log.Error(err, "update pipeline state failed")
			return errors.Convert(err)
		}
	}
	endTime := time.Now()
	log.Info("pipeline finished in %d ms: %v", endTime.UnixMilli()-startTime.UnixMilli(), err)
	return errors.Convert(err)
}
