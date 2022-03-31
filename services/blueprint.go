package services

import (
	"encoding/json"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm/clause"
)

type BlueprintQuery struct {
	Enable   bool `form:"enable"`
	Page     int  `form:"page"`
	PageSize int  `form:"pageSize"`
}

func CreateBlueprint(newBlueprint *models.InputBlueprint) (*models.Blueprint, error) {
	blueprint := models.Blueprint{
		Enable:     newBlueprint.Enable,
		CronConfig: newBlueprint.CronConfig,
		Name:       newBlueprint.Name,
	}
	var err error
	// update tasks state
	blueprint.Tasks, err = json.Marshal(newBlueprint.Tasks)
	if err != nil {
		return nil, err
	}

	err = db.Create(&blueprint).Error
	if err != nil {
		logger.Error("create pipline failed", err)
		return nil, errors.InternalError
	}
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return nil, errors.InternalError
	}

	return &blueprint, nil
}

func GetBlueprints(query *BlueprintQuery) ([]*models.Blueprint, int64, error) {
	blueprints := make([]*models.Blueprint, 0)
	db := db.Model(blueprints).Order("id DESC").Where("enable = ?", query.Enable)

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
	return blueprints, count, nil
}

func GetBlueprint(blueprintId uint64) (*models.Blueprint, error) {
	blueprint := &models.Blueprint{}
	err := db.Find(blueprint, blueprintId).Error
	if err != nil {
		return nil, err
	}
	return blueprint, nil
}

func ModifyBlueprint(newBlueprint *models.EditBlueprint) (*models.Blueprint, error) {
	blueprint := models.Blueprint{}
	err := db.Model(&models.Blueprint{}).
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
	cLog := logger.Global.Nested("blueprint")
	for _, e := range c.Entries() {
		c.Remove(e.ID)
	}
	c.Stop()
	for _, pp := range blueprints {
		var tasks [][]*models.NewTask
		err = json.Unmarshal(pp.Tasks, &tasks)
		if err != nil {
			cLog.Error("created cron job failed: %s", err)
			return err
		}
		_, err := c.AddFunc(pp.CronConfig, func() {
			newPipeline := models.NewPipeline{}
			newPipeline.Tasks = tasks
			newPipeline.Name = pp.Name
			newPipeline.BlueprintId = pp.ID
			pipeline, err := CreatePipeline(&newPipeline)
			// Return all created tasks to the User
			if err != nil {
				cLog.Error("created cron job failed: %s", err)
				return
			}
			err = RunPipeline(pipeline.ID)
			if err != nil {
				cLog.Error("run cron job failed: %s", err)
				return
			}
			cLog.Info("Run new cron job successfully")
		})
		if err != nil {
			cLog.Error("created cron job failed: %s", err)
			return err
		}
	}
	if len(blueprints) > 0 {
		c.Start()
	}
	return nil
}
