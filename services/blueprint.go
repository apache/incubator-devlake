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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type BlueprintQuery struct {
	Enable   *bool `form:"enable,omitempty"`
	Page     int   `form:"page"`
	PageSize int   `form:"pageSize"`
}

var blueprintLog = logger.Global.Nested("blueprint")
var vld = validator.New()

func CreateBlueprint(blueprint *models.Blueprint) error {
	err := validateBlueprint(blueprint)
	if err != nil {
		return err
	}
	err = db.Create(&blueprint).Error
	if err != nil {
		return err
	}
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return errors.InternalError
	}
	return nil
}

func GetBlueprints(query *BlueprintQuery) ([]*models.Blueprint, int64, error) {
	blueprints := make([]*models.Blueprint, 0)
	db := db.Model(blueprints).Order("id DESC")
	if query.Enable != nil {
		db = db.Where("enable = ?", *query.Enable)
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
	err = db.Find(&blueprints).Error
	if err != nil {
		return nil, 0, err
	}
	return blueprints, count, nil
}

func GetBlueprint(blueprintId uint64) (*models.Blueprint, error) {
	blueprint := &models.Blueprint{}
	err := db.First(blueprint, blueprintId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("blueprint not found")
		}
		return nil, err
	}
	return blueprint, nil
}

/*
func ModifyBlueprint(newBlueprint *models.EditBlueprint) (*models.Blueprint, error) {
	_, err := cron.ParseStandard(newBlueprint.CronConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid cronConfig: %w", err)
	}

	blueprint := models.Blueprint{}
	err = db.Model(&models.Blueprint{}).
		Where("id = ?", newBlueprint.BlueprintId).Limit(1).Find(&blueprint).Error
	if err != nil {
		return nil, err
	}
	// update cronConfig
	if newBlueprint.CronConfig != "" {
		blueprint.CronConfig = newBlueprint.CronConfig
	}
	// update tasks
	if newBlueprint.Tasks != nil {
		blueprint.Tasks, err = json.Marshal(newBlueprint.Tasks)
		if err != nil {
			return nil, err
		}
	}
	blueprint.Enable = newBlueprint.Enable

	err = db.Model(&models.Blueprint{}).
		Clauses(clause.OnConflict{UpdateAll: true}).Create(&blueprint).Error
	if err != nil {
		return nil, errors.InternalError
	}
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return nil, errors.InternalError
	}
	return &blueprint, nil
}
*/

func validateBlueprint(blueprint *models.Blueprint) error {
	// validation
	err := vld.Struct(blueprint)
	if err != nil {
		return err
	}
	_, err = cron.ParseStandard(blueprint.CronConfig)
	if err != nil {
		return fmt.Errorf("invalid cronConfig: %w", err)
	}
	tasks := make([][]models.NewTask, 0)
	err = json.Unmarshal(blueprint.Tasks, &tasks)
	if err != nil {
		return fmt.Errorf("invalid tasks: %w", err)
	}
	// tasks should not be empty
	if len(tasks) == 0 || len(tasks[0]) == 0 {
		return fmt.Errorf("empty tasks")
	}
	// TODO: validate each of every task object
	return nil
}

func UpdateBlueprint(blueprint *models.Blueprint) error {
	// validation
	err := validateBlueprint(blueprint)
	if err != nil {
		return err
	}
	// save
	err = db.Save(blueprint).Error
	if err != nil {
		return errors.InternalError
	}
	// reload schedule
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return errors.InternalError
	}
	// done
	return nil
}

func DeleteBlueprint(id uint64) error {
	err := db.Delete(&models.Blueprint{}, "id = ?", id).Error
	if err != nil {
		return errors.InternalError
	}
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return errors.InternalError
	}
	return nil
}

func ReloadBlueprints(c *cron.Cron) error {
	blueprints := make([]*models.Blueprint, 0)
	err := db.Model(&models.Blueprint{}).Where("enable = ?", true).Find(&blueprints).Error
	if err != nil {
		panic(err)
	}
	for _, e := range c.Entries() {
		c.Remove(e.ID)
	}
	c.Stop()
	for _, pp := range blueprints {
		var tasks [][]*models.NewTask
		err = json.Unmarshal(pp.Tasks, &tasks)
		if err != nil {
			blueprintLog.Error("created cron job failed: %s", err)
			return err
		}
		blueprint := pp
		_, err := c.AddFunc(pp.CronConfig, func() {
			newPipeline := models.NewPipeline{}
			newPipeline.Tasks = tasks
			newPipeline.Name = blueprint.Name
			newPipeline.BlueprintId = blueprint.ID
			pipeline, err := CreatePipeline(&newPipeline)
			// Return all created tasks to the User
			if err != nil {
				blueprintLog.Error("created cron job failed: %s", err)
				return
			}
			err = RunPipeline(pipeline.ID)
			if err != nil {
				blueprintLog.Error("run cron job failed: %s", err)
				return
			}
			blueprintLog.Info("Run new cron job successfully")
		})
		if err != nil {
			blueprintLog.Error("created cron job failed: %s", err)
			return err
		}
	}
	if len(blueprints) > 0 {
		c.Start()
	}
	return nil
}
