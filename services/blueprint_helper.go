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

// CreateDbBlueprint accepts a Blueprint instance and insert it to database
func CreateDbBlueprint(dbBlueprint *models.DbBlueprint) errors.Error {
	err := db.Create(&dbBlueprint).Error
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB blueprint")
	}
	return nil
}

// GetDbBlueprints returns a paginated list of Blueprints based on `query`
func GetDbBlueprints(query *BlueprintQuery) ([]*models.DbBlueprint, int64, errors.Error) {
	dbBlueprints := make([]*models.DbBlueprint, 0)
	db := db.Model(dbBlueprints).Order("id DESC")
	if query.Enable != nil {
		db = db.Where("enable = ?", *query.Enable)
	}

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of blueprints")
	}

	db = processDbClausesWithPager(db, query.PageSize, query.Page)

	err = db.Find(&dbBlueprints).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error finding DB blueprints")
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
	return dbBlueprint, nil
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
func parseBlueprint(DbBlueprint *models.DbBlueprint) *models.Blueprint {
	blueprint := models.Blueprint{
		Name:       DbBlueprint.Name,
		Mode:       DbBlueprint.Mode,
		Plan:       []byte(DbBlueprint.Plan),
		Enable:     DbBlueprint.Enable,
		CronConfig: DbBlueprint.CronConfig,
		IsManual:   DbBlueprint.IsManual,
		SkipOnFail: DbBlueprint.SkipOnFail,
		Settings:   []byte(DbBlueprint.Settings),
		Model:      DbBlueprint.Model,
	}
	return &blueprint
}

// parseDbBlueprint
func parseDbBlueprint(blueprint *models.Blueprint) *models.DbBlueprint {
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
