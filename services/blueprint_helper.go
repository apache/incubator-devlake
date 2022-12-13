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
	"fmt"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

// SaveDbBlueprint accepts a Blueprint instance and upsert it to database
func SaveDbBlueprint(dbBlueprint *models.DbBlueprint) errors.Error {
	var err error
	if dbBlueprint.ID != 0 {
		err = db.Save(&dbBlueprint).Error
	} else {
		err = db.Create(&dbBlueprint).Error
	}
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB blueprint")
	}
	err = db.Delete(&models.DbBlueprintLabel{}, `blueprint_id = ?`, dbBlueprint.ID).Error
	if err != nil {
		return errors.Default.Wrap(err, "error delete DB blueprint's old labelModels")
	}
	if len(dbBlueprint.Labels) > 0 {
		for i := range dbBlueprint.Labels {
			dbBlueprint.Labels[i].BlueprintId = dbBlueprint.ID
		}
		err = db.Create(&dbBlueprint.Labels).Error
		if err != nil {
			return errors.Default.Wrap(err, "error creating DB blueprint's labelModels")
		}
	}
	return nil
}

// GetDbBlueprints returns a paginated list of Blueprints based on `query`
func GetDbBlueprints(query *BlueprintQuery) ([]*models.DbBlueprint, int64, errors.Error) {
	dbBlueprints := make([]*models.DbBlueprint, 0)
	dbQuery := db.Model(dbBlueprints).Order("id DESC")
	if query.Enable != nil {
		dbQuery = dbQuery.Where("enable = ?", *query.Enable)
	}
	if query.IsManual != nil {
		dbQuery = dbQuery.Where("is_manual = ?", *query.IsManual)
	}
	if query.Label != "" {
		dbQuery = dbQuery.
			Joins(`left join _devlake_blueprint_labels ON _devlake_blueprint_labels.blueprint_id = _devlake_blueprints.id`).
			Where(`_devlake_blueprint_labels.name = ?`, query.Label)
	}

	var count int64
	err := dbQuery.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of blueprints")
	}

	dbQuery = processDbClausesWithPager(dbQuery, query.PageSize, query.Page)

	err = dbQuery.Find(&dbBlueprints).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error finding DB blueprints")
	}

	var blueprintIds []uint64
	for _, dbBlueprint := range dbBlueprints {
		blueprintIds = append(blueprintIds, dbBlueprint.ID)
	}
	var dbLabels []models.DbBlueprintLabel
	dbLabelsMap := map[uint64][]models.DbBlueprintLabel{}
	db.Where(`blueprint_id in ?`, blueprintIds).Find(&dbLabels)
	for _, dbLabel := range dbLabels {
		dbLabelsMap[dbLabel.BlueprintId] = append(dbLabelsMap[dbLabel.BlueprintId], dbLabel)
	}
	for _, dbBlueprint := range dbBlueprints {
		dbBlueprint.Labels = dbLabelsMap[dbBlueprint.ID]
	}

	return dbBlueprints, count, nil
}

// GetDbBlueprint returns the detail of a given Blueprint ID
func GetDbBlueprint(dbBlueprintId uint64) (*models.DbBlueprint, errors.Error) {
	dbBlueprint := &models.DbBlueprint{}
	err := db.First(dbBlueprint, dbBlueprintId).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.Wrap(err, "could not find blueprint in DB")
		}
		return nil, errors.Default.Wrap(err, "error getting blueprint from DB")
	}
	err = db.Find(&dbBlueprint.Labels, "blueprint_id = ?", dbBlueprint.ID).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting the blueprint labels from database")
	}
	return dbBlueprint, nil
}

// GetDbBlueprintByProjectName returns the detail of a given projectName
func GetDbBlueprintByProjectName(projectName string) (*models.DbBlueprint, errors.Error) {
	dbBlueprint := &models.DbBlueprint{}
	err := db.Where("project_name = ?", projectName).First(dbBlueprint).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find blueprint in DB by projectName %s", projectName))
		}
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error getting blueprint from DB by projectName %s", projectName))
	}
	err = db.Find(&dbBlueprint.Labels, "blueprint_id = ?", dbBlueprint.ID).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting the blueprint labels from database")
	}
	return dbBlueprint, nil
}

// RenameProjectNameForBlueprint FIXME ...
func RenameProjectNameForBlueprint(oldProjectName string, newProjectName string) errors.Error {
	err := db.Model(&models.DbBlueprint{}).
		Where("project_name = ?", oldProjectName).
		Updates(map[string]interface{}{
			"project_name": newProjectName,
		}).Error
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("Failed to RenameProjectNameForBlueprint from [%s] to [%s]", oldProjectName, newProjectName))
	}

	return nil
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
func parseBlueprint(dbBlueprint *models.DbBlueprint) *models.Blueprint {
	labelList := []string{}
	for _, labelModel := range dbBlueprint.Labels {
		labelList = append(labelList, labelModel.Name)
	}
	blueprint := models.Blueprint{
		Name:        dbBlueprint.Name,
		ProjectName: dbBlueprint.ProjectName,
		Mode:        dbBlueprint.Mode,
		Plan:        []byte(dbBlueprint.Plan),
		Enable:      dbBlueprint.Enable,
		CronConfig:  dbBlueprint.CronConfig,
		IsManual:    dbBlueprint.IsManual,
		SkipOnFail:  dbBlueprint.SkipOnFail,
		Settings:    []byte(dbBlueprint.Settings),
		Model:       dbBlueprint.Model,
		Labels:      labelList,
	}
	return &blueprint
}

// parseDbBlueprint
func parseDbBlueprint(blueprint *models.Blueprint) *models.DbBlueprint {
	dbBlueprint := models.DbBlueprint{
		Name:        blueprint.Name,
		ProjectName: blueprint.ProjectName,
		Mode:        blueprint.Mode,
		Plan:        string(blueprint.Plan),
		Enable:      blueprint.Enable,
		CronConfig:  blueprint.CronConfig,
		IsManual:    blueprint.IsManual,
		SkipOnFail:  blueprint.SkipOnFail,
		Settings:    string(blueprint.Settings),
		Model:       blueprint.Model,
	}
	dbBlueprint.Labels = []models.DbBlueprintLabel{}
	for _, label := range blueprint.Labels {
		dbBlueprint.Labels = append(dbBlueprint.Labels, models.DbBlueprintLabel{
			// NOTICE: BlueprintId may be nil
			BlueprintId: blueprint.ID,
			Name:        label,
		})
	}
	return &dbBlueprint
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
