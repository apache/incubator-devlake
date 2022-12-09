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
	"fmt"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

// ErrBlueprintRunning indicates there is a running pipeline with the specified blueprint_id
var ErrBlueprintRunning = errors.Default.New("the blueprint is running")

// CreateDbPipeline returns a NewPipeline
func CreateDbPipeline(newPipeline *models.NewPipeline) (*models.DbPipeline, errors.Error) {
	cronLocker.Lock()
	defer cronLocker.Unlock()
	if newPipeline.BlueprintId > 0 {
		var count int64
		err := db.Model(&models.DbPipeline{}).Where("blueprint_id = ? AND status IN ?", newPipeline.BlueprintId, models.PendingTaskStatus).Count(&count).Error
		if err != nil {
			return nil, errors.Default.Wrap(err, "query pipelines error")
		}
		if count > 0 {
			return nil, ErrBlueprintRunning
		}
	}
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
		globalPipelineLog.Error(err, "create pipeline failed: %v", err)
		return nil, errors.Internal.Wrap(err, "create pipeline failed")
	}

	dbPipeline.Labels = []models.DbPipelineLabel{}
	for _, label := range newPipeline.Labels {
		dbPipeline.Labels = append(dbPipeline.Labels, models.DbPipelineLabel{
			PipelineId: dbPipeline.ID,
			Name:       label,
		})
	}
	if len(dbPipeline.Labels) > 0 {
		if err := db.Create(&dbPipeline.Labels).Error; err != nil {
			globalPipelineLog.Error(err, "create pipeline's labelModels failed: %v", err)
			return nil, errors.Internal.Wrap(err, "create pipeline's labelModels failed")
		}
	}

	// create tasks accordingly
	for i := range newPipeline.Plan {
		for j := range newPipeline.Plan[i] {
			log.Debug(fmt.Sprintf("plan[%d][%d] is %+v\n", i, j, newPipeline.Plan[i][j]))
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
	dbQuery := db.Model(dbPipelines).Order("id DESC")
	if query.BlueprintId != 0 {
		dbQuery = dbQuery.Where("blueprint_id = ?", query.BlueprintId)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.Pending > 0 {
		dbQuery = dbQuery.Where("finished_at is null and status IN ?", models.PendingTaskStatus)
	}
	if query.Label != "" {
		dbQuery = dbQuery.
			Joins(`left join _devlake_pipeline_labels ON _devlake_pipeline_labels.pipeline_id = _devlake_pipelines.id`).
			Where(`_devlake_pipeline_labels.name = ?`, query.Label)
	}
	var count int64
	err := dbQuery.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB pipelines count")
	}

	dbQuery = processDbClausesWithPager(dbQuery, query.PageSize, query.Page)

	err = dbQuery.Find(&dbPipelines).Error
	if err != nil {
		return nil, count, errors.Default.Wrap(err, "error finding DB pipelines")
	}

	var pipelineIds []uint64
	for _, dbPipeline := range dbPipelines {
		pipelineIds = append(pipelineIds, dbPipeline.ID)
	}
	dbLabels := []models.DbPipelineLabel{}
	db.Where(`pipeline_id in ?`, pipelineIds).Find(&dbLabels)
	dbLabelsMap := map[uint64][]models.DbPipelineLabel{}
	for _, dbLabel := range dbLabels {
		dbLabelsMap[dbLabel.PipelineId] = append(dbLabelsMap[dbLabel.PipelineId], dbLabel)
	}
	for _, dbPipeline := range dbPipelines {
		dbPipeline.Labels = dbLabelsMap[dbPipeline.ID]
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
	err = db.Find(&dbPipeline.Labels, "pipeline_id = ?", pipelineId).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting the pipeline from database")
	}
	return dbPipeline, nil
}

// parsePipeline converts DbPipeline to Pipeline
func parsePipeline(dbPipeline *models.DbPipeline) *models.Pipeline {
	labelList := []string{}
	for _, labelModel := range dbPipeline.Labels {
		labelList = append(labelList, labelModel.Name)
	}
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
		Labels:        labelList,
	}
	return &pipeline
}

// parseDbPipeline converts Pipeline to DbPipeline
// nolint:unused
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
	dbPipeline.Labels = []models.DbPipelineLabel{}
	for _, label := range pipeline.Labels {
		dbPipeline.Labels = append(dbPipeline.Labels, models.DbPipelineLabel{
			// NOTICE: PipelineId may be nil
			PipelineId: pipeline.ID,
			Name:       label,
		})
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

// UpdateDbPipelineStatus update the status of pipeline
func UpdateDbPipelineStatus(pipelineId uint64, status string) errors.Error {
	err := db.Model(&models.DbPipeline{}).Where("id = ?", pipelineId).Update("status", status).Error
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}
