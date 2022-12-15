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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// RunPipeline FIXME ...
func RunPipeline(
	basicRes core.BasicRes,
	pipelineId uint64,
	runTasks func([]uint64) errors.Error,
) errors.Error {
	// load tasks for pipeline
	db := basicRes.GetDal()
	var tasks []models.Task
	err := db.All(&tasks, dal.Where("pipeline_id = ? AND status = ?", pipelineId, models.TASK_CREATED))
	if err != nil {
		return err
	}
	taskIds := make([][]uint64, 0)
	for _, task := range tasks {
		for len(taskIds) < task.PipelineRow {
			taskIds = append(taskIds, make([]uint64, 0))
		}
		taskIds[task.PipelineRow-1] = append(taskIds[task.PipelineRow-1], task.ID)
	}
	return runPipelineTasks(basicRes, pipelineId, taskIds, runTasks)
}

func runPipelineTasks(
	basicRes core.BasicRes,
	pipelineId uint64,
	taskIds [][]uint64,
	runTasks func([]uint64) errors.Error,
) errors.Error {
	db := basicRes.GetDal()
	log := basicRes.GetLogger()
	// load pipeline from db
	dbPipeline := &models.DbPipeline{}
	err := db.First(dbPipeline, dal.Where("id = ?", pipelineId))
	if err != nil {
		return err
	}

	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	for i, row := range taskIds {
		// update stage
		err = db.UpdateColumns(dbPipeline, []dal.DalSet{
			{ColumnName: "status", Value: models.TASK_RUNNING},
			{ColumnName: "stage", Value: i + 1},
		})
		if err != nil {
			log.Error(err, "update pipeline state failed")
			break
		}
		// run tasks in parallel
		err = runTasks(row)
		if err != nil {
			log.Error(err, "run tasks failed")
			return err
		}

		// update finishedTasks
		err = db.UpdateColumn(dbPipeline, "finished_tasks", dal.Expr("finished_tasks + ?", len(row)))
		if err != nil {
			log.Error(err, "update pipeline state failed")
			return err
		}
	}
	log.Info("pipeline finished in %d ms: %v", time.Now().UnixMilli()-dbPipeline.BeganAt.UnixMilli(), err)
	return err
}
