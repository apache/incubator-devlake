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

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// SaveDbBlueprint accepts a Blueprint instance and upsert it to database
func SaveDbBlueprint(dbBlueprint *models.DbBlueprint) errors.Error {
	var err error
	if dbBlueprint.ID != 0 {
		err = db.Update(&dbBlueprint)
	} else {
		err = db.Create(&dbBlueprint)
	}
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB blueprint")
	}
	err = db.Delete(&models.DbBlueprintLabel{}, dal.Where(`blueprint_id = ?`, dbBlueprint.ID))
	if err != nil {
		return errors.Default.Wrap(err, "error delete DB blueprint's old labelModels")
	}
	if len(dbBlueprint.Labels) > 0 {
		for i := range dbBlueprint.Labels {
			dbBlueprint.Labels[i].BlueprintId = dbBlueprint.ID
		}
		err = db.Create(&dbBlueprint.Labels)
		if err != nil {
			return errors.Default.Wrap(err, "error creating DB blueprint's labelModels")
		}
	}
	return nil
}

// GetDbBlueprints returns a paginated list of Blueprints based on `query`
func GetDbBlueprints(query *BlueprintQuery) ([]*models.DbBlueprint, int64, errors.Error) {
	// process query parameters
	clauses := []dal.Clause{dal.From(&models.DbBlueprint{})}
	if query.Enable != nil {
		clauses = append(clauses, dal.Where("enable = ?", *query.Enable))
	}
	if query.IsManual != nil {
		clauses = append(clauses, dal.Where("is_manual = ?", *query.IsManual))
	}
	if query.Label != "" {
		clauses = append(clauses,
			dal.Join("left join _devlake_blueprint_labels bl ON bl.blueprint_id = _devlake_blueprints.id"),
			dal.Where("bl.name = ?", query.Label),
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
	dbBlueprints := make([]*models.DbBlueprint, 0)
	err = db.All(&dbBlueprints, clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of blueprints")
	}

	// load labels for blueprints
	for _, dbBlueprint := range dbBlueprints {
		err = fillBlueprintDetail(dbBlueprint)
		if err != nil {
			return nil, 0, err
		}
	}

	return dbBlueprints, count, nil
}

// GetDbBlueprint returns the detail of a given Blueprint ID
func GetDbBlueprint(dbBlueprintId uint64) (*models.DbBlueprint, errors.Error) {
	dbBlueprint := &models.DbBlueprint{}
	err := db.First(dbBlueprint, dal.Where("id = ?", dbBlueprintId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, "could not find blueprint in DB")
		}
		return nil, errors.Default.Wrap(err, "error getting blueprint from DB")
	}
	err = fillBlueprintDetail(dbBlueprint)
	if err != nil {
		return nil, err
	}
	return dbBlueprint, nil
}

// GetDbBlueprintByProjectName returns the detail of a given projectName
func GetDbBlueprintByProjectName(projectName string) (*models.DbBlueprint, errors.Error) {
	dbBlueprint := &models.DbBlueprint{}
	err := db.First(dbBlueprint, dal.Where("project_name = ?", projectName))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find blueprint in DB by projectName %s", projectName))
		}
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error getting blueprint from DB by projectName %s", projectName))
	}
	err = fillBlueprintDetail(dbBlueprint)
	if err != nil {
		return nil, err
	}
	return dbBlueprint, nil
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
	if len(blueprint.Settings) == 0 {
		blueprint.Settings = nil
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

func fillBlueprintDetail(blueprint *models.DbBlueprint) errors.Error {
	err := db.All(&blueprint.Labels, dal.Where("blueprint_id = ?", blueprint.ID))
	if err != nil {
		return errors.Internal.Wrap(err, "error getting the blueprint labels from database")
	}
	return nil
}
