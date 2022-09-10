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
	"fmt"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

// CreateDbPipeline returns a NewPipeline
func CreateDbPipeline(newPipeline *models.NewPipeline) (*models.DbPipeline, error) {
	// create pipeline object from posted data
	dbPipeline := &models.DbPipeline{
		Name:          newPipeline.Name,
		FinishedTasks: 0,
		Status:        models.TASK_CREATED,
		Message:       "",
		SpentSeconds:  0,
	}
	if newPipeline.BlueprintId != 0 {
		dbPipeline.BlueprintId = newPipeline.BlueprintId
	}
	dbPipeline, err := encryptDbPipeline(dbPipeline)
	if err != nil {
		return nil, err
	}
	// save pipeline to database
	err = db.Create(&dbPipeline).Error
	if err != nil {
		globalPipelineLog.Error(err, "create pipline failed: %w", err)
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
				globalPipelineLog.Error(err, "create task for pipeline failed: %w", err)
				return nil, err
			}
			// sync task state back to pipeline
			dbPipeline.TotalTasks += 1
		}
	}
	if err != nil {
		globalPipelineLog.Error(err, "save tasks for pipeline failed: %w", err)
		return nil, errors.Internal.Wrap(err, "save tasks for pipeline failed")
	}
	if dbPipeline.TotalTasks == 0 {
		return nil, fmt.Errorf("no task to run")
	}

	// update tasks state
	planByte, err := json.Marshal(newPipeline.Plan)
	if err != nil {
		return nil, err
	}
	dbPipeline.Plan = string(planByte)
	dbPipeline, err = encryptDbPipeline(dbPipeline)
	if err != nil {
		return nil, err
	}
	err = db.Model(dbPipeline).Updates(map[string]interface{}{
		"total_tasks": dbPipeline.TotalTasks,
		"plan":        dbPipeline.Plan,
	}).Error
	if err != nil {
		globalPipelineLog.Error(err, "update pipline state failed: %w", err)
		return nil, errors.Internal.Wrap(err, "update pipline state failed")
	}

	return dbPipeline, nil
}

// GetDbPipelines by query
func GetDbPipelines(query *PipelineQuery) ([]*models.DbPipeline, int64, error) {
	dbPipelines := make([]*models.DbPipeline, 0)
	db := db.Model(dbPipelines).Order("id DESC")
	if query.BlueprintId != 0 {
		db = db.Where("blueprint_id = ?", query.BlueprintId)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Pending > 0 {
		db = db.Where("finished_at is null")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if query.Page > 0 && query.PageSize > 0 {
		offset := query.PageSize * (query.Page - 1)
		db = db.Limit(query.PageSize).Offset(offset)
	}
	err = db.Find(&dbPipelines).Error
	if err != nil {
		return nil, count, err
	}
	return dbPipelines, count, nil
}

// GetDbPipeline by id
func GetDbPipeline(pipelineId uint64) (*models.DbPipeline, error) {
	dbPipeline := &models.DbPipeline{}
	err := db.First(dbPipeline, pipelineId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound.New("pipeline not found", errors.AsUserMessage())
		}
		return nil, errors.Internal.Wrap(err, "error getting the pipeline from database", errors.AsUserMessage())
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
func encryptDbPipeline(dbPipeline *models.DbPipeline) (*models.DbPipeline, error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	planEncrypt, err := core.Encrypt(encKey, dbPipeline.Plan)
	if err != nil {
		return nil, err
	}
	dbPipeline.Plan = planEncrypt
	return dbPipeline, nil
}

// encryptDbPipeline decrypts dbPipeline.Plan
func decryptDbPipeline(dbPipeline *models.DbPipeline) (*models.DbPipeline, error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	plan, err := core.Decrypt(encKey, dbPipeline.Plan)
	if err != nil {
		return nil, err
	}
	dbPipeline.Plan = plan
	return dbPipeline, nil
}
