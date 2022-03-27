package services

import (
	"encoding/json"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"gorm.io/gorm/clause"
	"strings"
)

func init() {
	v := config.GetConfig()
	var notificationEndpoint = v.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = v.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}
	db.Model(&models.Pipeline{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)
}

func CreatePipelinePlan(newPipeline *models.NewPipeline) (*models.PipelinePlan, error) {
	pipelinePlan := &models.PipelinePlan{
		Enable:     newPipeline.Enable,
		CronConfig: newPipeline.CronConfig,
		Name:       newPipeline.Name,
	}
	var err error
	// update tasks state
	pipelinePlan.Tasks, err = json.Marshal(newPipeline.Tasks)
	if err != nil {
		return nil, err
	}

	err = db.Create(&pipelinePlan).Error
	if err != nil {
		logger.Error("create pipline failed", err)
		return nil, errors.InternalError
	}
	return pipelinePlan, nil
}

func GetPipelinePlans() ([]*models.PipelinePlan, int64, error) {
	pipelinePlans := make([]*models.PipelinePlan, 0)
	db := db.Model(pipelinePlans).Order("id DESC")

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&pipelinePlans).Error
	if err != nil {
		return nil, count, err
	}
	return pipelinePlans, count, nil
}

func GetPipelinePlan(pipelinePlanId uint64) (*models.PipelinePlan, error) {
	pipelinePlan := &models.PipelinePlan{}
	err := db.Find(pipelinePlan, pipelinePlanId).Error
	if err != nil {
		return nil, err
	}
	return pipelinePlan, nil
}

func ModifyPipelinePlan(newPipeline *models.NewPipeline, pipelinePlanId uint64) (*models.PipelinePlan, error) {
	pipelinePlan := &models.PipelinePlan{}
	err := db.Model(&models.PipelinePlan{}).
		Where("id = ?", pipelinePlanId).Limit(1).Find(pipelinePlan).Error
	if err != nil {
		return nil, err
	}
	pipelinePlan.CronConfig = newPipeline.CronConfig
	pipelinePlan.Enable = newPipeline.Enable
	// update tasks state
	pipelinePlan.Tasks, err = json.Marshal(newPipeline.Tasks)
	if err != nil {
		return nil, err
	}

	err = db.Model(&models.PipelinePlan{}).
		Clauses(clause.OnConflict{UpdateAll: true}).Create(&pipelinePlan).Error
	if err != nil {
		logger.Error("create pipline failed", err)
		return nil, errors.InternalError
	}
	return pipelinePlan, nil
}

func DeletePipelinePlan(id uint64) error {
	err := db.Delete(&models.PipelinePlan{}, "id = ?", id).Error
	if err != nil {
		logger.Error("create pipline failed", err)
		return errors.InternalError
	}
	return nil
}
