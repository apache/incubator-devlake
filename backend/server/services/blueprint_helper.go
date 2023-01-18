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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
)

// SaveDbBlueprint accepts a Blueprint instance and upsert it to database
func SaveDbBlueprint(blueprint *models.Blueprint) errors.Error {
	var err error
	if blueprint.ID != 0 {
		err = db.Update(&blueprint)
	} else {
		err = db.Create(&blueprint)
	}
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB blueprint")
	}
	err = db.Delete(&models.DbBlueprintLabel{}, dal.Where(`blueprint_id = ?`, blueprint.ID))
	if err != nil {
		return errors.Default.Wrap(err, "error delete DB blueprint's old labelModels")
	}
	if len(blueprint.Labels) > 0 {
		blueprintLabels := make([]*models.DbBlueprintLabel, 0)
		for i := range blueprint.Labels {
			blueprintLabels = append(blueprintLabels, &models.DbBlueprintLabel{
				BlueprintId: blueprint.ID,
				Name:        blueprint.Labels[i],
			})
		}
		err = db.Create(&blueprintLabels)
		if err != nil {
			return errors.Default.Wrap(err, "error creating DB blueprint's labelModels")
		}
	}
	return nil
}

// GetDbBlueprints returns a paginated list of Blueprints based on `query`
func GetDbBlueprints(query *BlueprintQuery) ([]*models.Blueprint, int64, errors.Error) {
	// process query parameters
	clauses := []dal.Clause{dal.From(&models.Blueprint{})}
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
	dbBlueprints := make([]*models.Blueprint, 0)
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
func GetDbBlueprint(blueprintId uint64) (*models.Blueprint, errors.Error) {
	blueprint := &models.Blueprint{}
	err := db.First(blueprint, dal.Where("id = ?", blueprintId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, "could not find blueprint in DB")
		}
		return nil, errors.Default.Wrap(err, "error getting blueprint from DB")
	}
	err = fillBlueprintDetail(blueprint)
	if err != nil {
		return nil, err
	}
	return blueprint, nil
}

// GetDbBlueprintByProjectName returns the detail of a given projectName
func GetDbBlueprintByProjectName(projectName string) (*models.Blueprint, errors.Error) {
	dbBlueprint := &models.Blueprint{}
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

func fillBlueprintDetail(blueprint *models.Blueprint) errors.Error {
	err := db.Pluck("name", &blueprint.Labels, dal.From(&models.DbBlueprintLabel{}), dal.Where("blueprint_id = ?", blueprint.ID))
	if err != nil {
		return errors.Internal.Wrap(err, "error getting the blueprint labels from database")
	}
	return nil
}
