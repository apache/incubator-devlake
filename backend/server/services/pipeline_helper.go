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
	"fmt"
	"sync"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/helpers/dbhelper"
)

var createDbPipelineLock sync.Mutex

// CreateDbPipeline returns a NewPipeline
func CreateDbPipeline(newPipeline *models.NewPipeline) (pipeline *models.Pipeline, err errors.Error) {
	createDbPipelineLock.Lock()
	defer createDbPipelineLock.Unlock()
	pipeline = &models.Pipeline{}
	txHelper := dbhelper.NewTxHelper(basicRes, &err)
	defer txHelper.End()
	tx := txHelper.Begin()
	errors.Must(txHelper.LockTablesTimeout(2*time.Second,
		dal.LockTables{
			{Table: "_devlake_pipelines", Exclusive: true},
			{Table: "_devlake_pipeline_labels", Exclusive: true},
			{Table: "_devlake_tasks", Exclusive: true},
		},
	))
	if err != nil {
		err = errors.BadInput.Wrap(err, "failed to lock pipeline table, is there any pending pipeline or deletion?")
		return
	}
	if newPipeline.BlueprintId > 0 {
		count := errors.Must1(tx.Count(
			dal.From(&models.Pipeline{}),
			dal.Where("blueprint_id = ? AND status IN ?", newPipeline.BlueprintId, models.PendingTaskStatus),
		))
		// some pipeline is running, get the detail and output them.
		if count > 0 {
			return nil, errors.BadInput.New("there are pending pipelines of current blueprint already")
		}
	}
	// create pipeline object from posted data
	dbPipeline := &models.Pipeline{
		Name:          newPipeline.Name,
		FinishedTasks: 0,
		Status:        models.TASK_CREATED,
		Message:       "",
		SpentSeconds:  0,
		Plan:          newPipeline.Plan,
		SyncPolicy:    newPipeline.SyncPolicy,
	}
	if newPipeline.BlueprintId != 0 {
		dbPipeline.BlueprintId = newPipeline.BlueprintId
	}

	// save pipeline to database
	errors.Must(tx.Create(dbPipeline))
	labels := make([]models.DbPipelineLabel, 0)
	for _, label := range newPipeline.Labels {
		labels = append(labels, models.DbPipelineLabel{
			PipelineId: dbPipeline.ID,
			Name:       label,
		})
	}
	if len(newPipeline.Labels) > 0 {
		errors.Must(tx.Create(&labels))
	}

	// create tasks accordingly
	for i := range newPipeline.Plan {
		for j := range newPipeline.Plan[i] {
			logger.Debug(fmt.Sprintf("plan[%d][%d] is %+v\n", i, j, newPipeline.Plan[i][j]))
			pipelineTask := newPipeline.Plan[i][j]
			newTask := &models.NewTask{
				PipelineTask: pipelineTask,
				PipelineId:   dbPipeline.ID,
				PipelineRow:  i + 1,
				PipelineCol:  j + 1,
			}
			_ = errors.Must1(createTask(newTask, tx))
			// sync task state back to pipeline
			dbPipeline.TotalTasks += 1
		}
	}
	if dbPipeline.TotalTasks == 0 {
		err = errors.BadInput.New("no task to run")
		return
	}

	// update tasks state
	errors.Must(tx.Update(dbPipeline))
	dbPipeline.Labels = newPipeline.Labels
	return dbPipeline, nil
}

// GetDbPipelines by query
func GetDbPipelines(query *PipelineQuery) ([]*models.Pipeline, int64, errors.Error) {
	// process query parameters
	clauses := []dal.Clause{dal.From(&models.Pipeline{})}
	if query.BlueprintId != 0 {
		clauses = append(clauses, dal.Where("blueprint_id = ?", query.BlueprintId))
	}
	if query.Status != "" {
		clauses = append(clauses, dal.Where("status = ?", query.Status))
	}
	if query.Pending > 0 {
		clauses = append(clauses, dal.Where("finished_at is null and status IN ?", models.PendingTaskStatus))
	}
	if query.Label != "" {
		clauses = append(clauses,
			dal.Join("LEFT JOIN _devlake_pipeline_labels pl ON pl.pipeline_id = _devlake_pipelines.id"),
			dal.Where("pl.name = ?", query.Label),
		)
	}

	// count total records
	count, err := db.Count(clauses...)
	if err != nil {
		return nil, 0, err
	}

	// load paginated blueprints from database
	clauses = append(clauses,
		dal.Orderby("id DESC"),
		dal.Offset(query.GetSkip()),
		dal.Limit(query.GetPageSize()),
	)
	dbPipelines := make([]*models.Pipeline, 0)
	err = db.All(&dbPipelines, clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of pipelines")
	}

	// load labels for blueprints
	for _, dbPipeline := range dbPipelines {
		err = fillPipelineDetail(dbPipeline)
		if err != nil {
			return nil, 0, err
		}
	}

	return dbPipelines, count, nil
}

// GetDbPipeline by id
func GetDbPipeline(pipelineId uint64) (*models.Pipeline, errors.Error) {
	dbPipeline := &models.Pipeline{}
	err := db.First(dbPipeline, dal.Where("id = ?", pipelineId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.New("pipeline not found")
		}
		return nil, errors.Internal.Wrap(err, "error getting the pipeline from database")
	}
	err = fillPipelineDetail(dbPipeline)
	if err != nil {
		return nil, err
	}
	return dbPipeline, nil
}

func fillPipelineDetail(pipeline *models.Pipeline) errors.Error {
	err := basicRes.GetDal().Pluck("name", &pipeline.Labels, dal.From(&models.DbPipelineLabel{}), dal.Where("pipeline_id = ?", pipeline.ID))
	if err != nil {
		return errors.Internal.Wrap(err, "error getting the pipeline labels from database")
	}
	return nil
}
