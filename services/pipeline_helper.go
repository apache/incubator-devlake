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
	"encoding/json"
	goerror "errors"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

// CreateDbPipeline returns a NewPipeline
func CreateDbPipeline(newPipeline *models.NewPipeline) (*models.DbPipeline, errors.Error) {
	cronLocker.Lock()
	defer cronLocker.Unlock()
	planByte, err := errors.Convert01(json.Marshal(newPipeline.Plan))
	if err != nil {
		return nil, err
	}
	// create pipeline object from posted data
	dbPipeline := &models.DbPipeline{
		Name:          newPipeline.Name,
		FinishedTasks: 0,
		Status:        models.TASK_CREATED,
		Message:       "",
		SpentSeconds:  0,
		Plan:          string(planByte),
	}
	if newPipeline.BlueprintId != 0 {
		dbPipeline.BlueprintId = newPipeline.BlueprintId
	}
	dbPipeline, err = encryptDbPipeline(dbPipeline)
	if err != nil {
		return nil, err
	}
	// save pipeline to database
	if err := db.Create(&dbPipeline).Error; err != nil {
		globalPipelineLog.Error(err, "create pipline failed: %v", err)
		return nil, errors.Internal.Wrap(err, "create pipline failed")
	}

	// create tasks accordingly
	for i := range newPipeline.Plan {
		for j := range newPipeline.Plan[i] {
			pipelineTask := newPipeline.Plan[i][j]
			newTask := &models.NewTask{
				PipelineTask: pipelineTask,
				PipelineId:   dbPipeline.ID,
				PipelineRow:  i + 1,
				PipelineCol:  j + 1,
			}
			_, err := CreateTask(newTask)
			if err != nil {
				globalPipelineLog.Error(err, "create task for pipeline failed: %v", err)
				return nil, err
			}
			// sync task state back to pipeline
			dbPipeline.TotalTasks += 1
		}
	}
	if err != nil {
		globalPipelineLog.Error(err, "save tasks for pipeline failed: %v", err)
		return nil, errors.Internal.Wrap(err, "save tasks for pipeline failed")
	}
	if dbPipeline.TotalTasks == 0 {
		return nil, errors.Internal.New("no task to run")
	}

	// update tasks state
	if err := db.Model(dbPipeline).Updates(map[string]interface{}{
		"total_tasks": dbPipeline.TotalTasks,
	}).Error; err != nil {
		globalPipelineLog.Error(err, "update pipline state failed: %v", err)
		return nil, errors.Internal.Wrap(err, "update pipline state failed")
	}

	return dbPipeline, nil
}

// GetDbPipelines by query
func GetDbPipelines(query *PipelineQuery) ([]*models.DbPipeline, int64, errors.Error) {
	dbPipelines := make([]*models.DbPipeline, 0)
	db := db.Model(dbPipelines).Order("id DESC")
	if query.BlueprintId != 0 {
		db = db.Where("blueprint_id = ?", query.BlueprintId)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Pending > 0 {
		db = db.Where("finished_at is null and status != TASK_FAILED")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB pipelines count")
	}

	db = processDbClausesWithPager(db, query.PageSize, query.Page)

	err = db.Find(&dbPipelines).Error
	if err != nil {
		return nil, count, errors.Default.Wrap(err, "error finding DB pipelines")
	}
	return dbPipelines, count, nil
}

// GetDbPipeline by id
func GetDbPipeline(pipelineId uint64) (*models.DbPipeline, errors.Error) {
	dbPipeline := &models.DbPipeline{}
	err := db.First(dbPipeline, pipelineId).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.New("pipeline not found")
		}
		return nil, errors.Internal.Wrap(err, "error getting the pipeline from database")
	}
	return dbPipeline, nil
}

// parsePipeline converts DbPipeline to Pipeline
func parsePipeline(dbPipeline *models.DbPipeline) *models.Pipeline {
	pipeline := models.Pipeline{
		Model:         dbPipeline.Model,
		Name:          dbPipeline.Name,
		BlueprintId:   dbPipeline.BlueprintId,
		Plan:          []byte(dbPipeline.Plan),
		TotalTasks:    dbPipeline.TotalTasks,
		FinishedTasks: dbPipeline.FinishedTasks,
		BeganAt:       dbPipeline.BeganAt,
		FinishedAt:    dbPipeline.FinishedAt,
		Status:        dbPipeline.Status,
		Message:       dbPipeline.Message,
		SpentSeconds:  dbPipeline.SpentSeconds,
		Stage:         dbPipeline.Stage,
	}
	return &pipeline
}

// parseDbPipeline converts Pipeline to DbPipeline
func parseDbPipeline(pipeline *models.Pipeline) *models.DbPipeline {
	dbPipeline := models.DbPipeline{
		Model:         pipeline.Model,
		Name:          pipeline.Name,
		BlueprintId:   pipeline.BlueprintId,
		Plan:          string(pipeline.Plan),
		TotalTasks:    pipeline.TotalTasks,
		FinishedTasks: pipeline.FinishedTasks,
		BeganAt:       pipeline.BeganAt,
		FinishedAt:    pipeline.FinishedAt,
		Status:        pipeline.Status,
		Message:       pipeline.Message,
		SpentSeconds:  pipeline.SpentSeconds,
		Stage:         pipeline.Stage,
	}
	return &dbPipeline
}

// encryptDbPipeline encrypts dbPipeline.Plan
func encryptDbPipeline(dbPipeline *models.DbPipeline) (*models.DbPipeline, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	planEncrypt, err := core.Encrypt(encKey, dbPipeline.Plan)
	if err != nil {
		return nil, err
	}
	dbPipeline.Plan = planEncrypt
	return dbPipeline, nil
}

// encryptDbPipeline decrypts dbPipeline.Plan
func decryptDbPipeline(dbPipeline *models.DbPipeline) (*models.DbPipeline, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	plan, err := core.Decrypt(encKey, dbPipeline.Plan)
	if err != nil {
		return nil, err
	}
	dbPipeline.Plan = plan
	return dbPipeline, nil
}
