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
	goerror "errors"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

// SaveDbBlueprint accepts a Blueprint instance and upsert it to database
func SaveDbBlueprint(dbBlueprint *models.DbBlueprint, parallelLabelModels []models.DbBlueprintParallelLabel) errors.Error {
	err := db.Save(&dbBlueprint).Error
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB blueprint")
	}
	err = db.Delete(&models.DbBlueprintParallelLabel{}, `blueprint_id = ?`, dbBlueprint.ID).Error
	if err != nil {
		return errors.Default.Wrap(err, "error delete DB blueprint's old parallelLabelModels")
	}
	if len(parallelLabelModels) > 0 {
		for _, parallelLabelModel := range parallelLabelModels {
			parallelLabelModel.BlueprintId = dbBlueprint.ID
		}
		err = db.Create(&parallelLabelModels).Error
		if err != nil {
			return errors.Default.Wrap(err, "error creating DB blueprint's parallelLabelModels")
		}
	}
	return nil
}

// GetDbBlueprints returns a paginated list of Blueprints based on `query`
func GetDbBlueprints(query *BlueprintQuery) ([]*models.DbBlueprint, map[uint64][]models.DbBlueprintParallelLabel, int64, errors.Error) {
	dbBlueprints := make([]*models.DbBlueprint, 0)
	dbQuery := db.Model(dbBlueprints).Order("id DESC")
	if query.Enable != nil {
		dbQuery = dbQuery.Where("enable = ?", *query.Enable)
	}
	if query.IsManual != nil {
		dbQuery = dbQuery.Where("is_manual = ?", *query.IsManual)
	}
	if query.ParallelLabel != "" {
		dbQuery = dbQuery.
			Joins(`left join _devlake_blueprint_parallel_labels ON
                  _devlake_blueprint_parallel_labels.blueprint_id = _devlake_blueprints.id`).
			Where(`_devlake_blueprint_parallel_labels.name = ?`, query.ParallelLabel)
	}

	var count int64
	err := dbQuery.Count(&count).Error
	if err != nil {
		return nil, nil, 0, errors.Default.Wrap(err, "error getting DB count of blueprints")
	}

	dbQuery = processDbClausesWithPager(dbQuery, query.PageSize, query.Page)

	err = dbQuery.Find(&dbBlueprints).Error
	if err != nil {
		return nil, nil, 0, errors.Default.Wrap(err, "error finding DB blueprints")
	}

	var blueprintIds []uint64
	for _, dbBlueprint := range dbBlueprints {
		blueprintIds = append(blueprintIds, dbBlueprint.ID)
	}
	dbParallelLabels := []models.DbBlueprintParallelLabel{}
	dbParallelLabelsMap := map[uint64][]models.DbBlueprintParallelLabel{}
	db.Where(`blueprint_id in ?`, blueprintIds).Find(&dbParallelLabels)
	for _, dbParallelLabel := range dbParallelLabels {
		dbParallelLabelsMap[dbParallelLabel.BlueprintId] = append(dbParallelLabelsMap[dbParallelLabel.BlueprintId], dbParallelLabel)
	}

	return dbBlueprints, dbParallelLabelsMap, count, nil
}

// GetDbBlueprint returns the detail of a given Blueprint ID
func GetDbBlueprint(dbBlueprintId uint64) (*models.DbBlueprint, []models.DbBlueprintParallelLabel, errors.Error) {
	dbBlueprint := &models.DbBlueprint{}
	err := db.First(dbBlueprint, dbBlueprintId).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.NotFound.Wrap(err, "could not find blueprint in DB")
		}
		return nil, nil, errors.Default.Wrap(err, "error getting blueprint from DB")
	}
	dbParallelLabels := []models.DbBlueprintParallelLabel{}
	err = db.Find(&dbParallelLabels, "blueprint_id = ?", dbBlueprint.ID).Error
	if err != nil {
		return nil, nil, errors.Internal.Wrap(err, "error getting the blueprint from database")
	}
	return dbBlueprint, dbParallelLabels, nil
}

// DeleteDbBlueprint deletes blueprint by id
func DeleteDbBlueprint(id uint64) errors.Error {
	err := db.Delete(&models.DbBlueprint{}, "id = ?", id).Error
	if err != nil {
		return errors.Default.Wrap(err, "error deleting blueprint from DB")
	}
	return nil
}

// parseBlueprint
func parseBlueprint(DbBlueprint *models.DbBlueprint, parallelLabelModels []models.DbBlueprintParallelLabel) *models.Blueprint {
	parallelLabelList := []string{}
	for _, labelModel := range parallelLabelModels {
		parallelLabelList = append(parallelLabelList, labelModel.Name)
	}
	blueprint := models.Blueprint{
		Name:           DbBlueprint.Name,
		Mode:           DbBlueprint.Mode,
		Plan:           []byte(DbBlueprint.Plan),
		Enable:         DbBlueprint.Enable,
		CronConfig:     DbBlueprint.CronConfig,
		IsManual:       DbBlueprint.IsManual,
		SkipOnFail:     DbBlueprint.SkipOnFail,
		Settings:       []byte(DbBlueprint.Settings),
		Model:          DbBlueprint.Model,
		ParallelLabels: parallelLabelList,
	}
	return &blueprint
}

// parseDbBlueprint
func parseDbBlueprint(blueprint *models.Blueprint) (*models.DbBlueprint, []models.DbBlueprintParallelLabel) {
	dbBlueprint := models.DbBlueprint{
		Name:       blueprint.Name,
		Mode:       blueprint.Mode,
		Plan:       string(blueprint.Plan),
		Enable:     blueprint.Enable,
		CronConfig: blueprint.CronConfig,
		IsManual:   blueprint.IsManual,
		SkipOnFail: blueprint.SkipOnFail,
		Settings:   string(blueprint.Settings),
		Model:      blueprint.Model,
	}
	parallelLabels := []models.DbBlueprintParallelLabel{}
	for _, label := range blueprint.ParallelLabels {
		parallelLabels = append(parallelLabels, models.DbBlueprintParallelLabel{
			// NOTICE: BlueprintId may be nil
			BlueprintId: blueprint.ID,
			Name:        label,
		})
	}
	return &dbBlueprint, parallelLabels
}

// encryptDbBlueprint
func encryptDbBlueprint(dbBlueprint *models.DbBlueprint) (*models.DbBlueprint, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	planEncrypt, err := core.Encrypt(encKey, dbBlueprint.Plan)
	if err != nil {
		return nil, err
	}
	dbBlueprint.Plan = planEncrypt
	settingsEncrypt, err := core.Encrypt(encKey, dbBlueprint.Settings)
	dbBlueprint.Settings = settingsEncrypt
	if err != nil {
		return nil, err
	}
	return dbBlueprint, nil
}

// decryptDbBlueprint
func decryptDbBlueprint(dbBlueprint *models.DbBlueprint) (*models.DbBlueprint, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)
	plan, err := core.Decrypt(encKey, dbBlueprint.Plan)
	if err != nil {
		return nil, err
	}
	dbBlueprint.Plan = plan
	settings, err := core.Decrypt(encKey, dbBlueprint.Settings)
	dbBlueprint.Settings = settings
	if err != nil {
		return nil, err
	}
	return dbBlueprint, nil
}
