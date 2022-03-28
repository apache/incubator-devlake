package services

import (
	"encoding/json"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm/clause"
)

type PipelinePlanQuery struct {
	Enable   bool `form:"enable"`
	Page     int  `form:"page"`
	PageSize int  `form:"pageSize"`
}

func CreatePipelinePlan(newPipeline *models.InputPipelinePlan) (*models.PipelinePlan, error) {
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
	err = ReloadPipelinePlans()
	if err != nil {
		logger.Error("create cron job failed", err)
		return nil, errors.InternalError
	}

	return pipelinePlan, nil
}

func GetPipelinePlans(query *PipelinePlanQuery) ([]*models.PipelinePlan, int64, error) {
	pipelinePlans := make([]*models.PipelinePlan, 0)
	db := db.Model(pipelinePlans).Order("id DESC").Where("enable = ?", query.Enable)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if query.Page > 0 && query.PageSize > 0 {
		offset := query.PageSize * (query.Page - 1)
		db = db.Limit(query.PageSize).Offset(offset)
	}
	err = db.Find(&pipelinePlans).Error
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

func ModifyPipelinePlan(newPipelinePlan *models.EditPipelinePlan) (*models.PipelinePlan, error) {
	pipelinePlan := &models.PipelinePlan{}
	err := db.Model(&models.PipelinePlan{}).
		Where("id = ?", newPipelinePlan.PipelinePlanId).Limit(1).Find(pipelinePlan).Error
	if err != nil {
		return nil, err
	}
	// update cronConfig
	if newPipelinePlan.CronConfig != "" {
		pipelinePlan.CronConfig = newPipelinePlan.CronConfig
	}
	// update tasks
	if newPipelinePlan.Tasks != nil {
		pipelinePlan.Tasks, err = json.Marshal(newPipelinePlan.Tasks)
		if err != nil {
			return nil, err
		}
	}
	pipelinePlan.Enable = newPipelinePlan.Enable

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

func ReloadPipelinePlans() error {
	pipelinePlans := make([]*models.PipelinePlan, 0)
	err := db.Model(&models.PipelinePlan{}).Where("enable = ?", true).Find(&pipelinePlans).Error
	if err != nil {
		panic(err)
	}
	c := cron.New()
	cLog := logger.Global.Nested("plan")

	for _, pp := range pipelinePlans {
		var tasks [][]*models.NewTask
		err = json.Unmarshal(pp.Tasks, &tasks)
		if err != nil {
			cLog.Error("created cron job failed: %s", err)
			return err
		}
		//
		newPipeline := models.NewPipeline{}
		newPipeline.Tasks = tasks
		newPipeline.Name = pp.Name
		newPipeline.PipelinePlanId = pp.ID
		_, err = c.AddFunc(pp.CronConfig, func() {
			pipeline, err := CreatePipeline(&newPipeline)
			// Return all created tasks to the User
			if err != nil {
				cLog.Error("created cron job failed: %s", err)
				return
			}
			go func() {
				_ = RunPipeline(pipeline.ID)
			}()
		})
		if err != nil {
			cLog.Error("created cron job failed: %s", err)
			return err

		}
	}
	if len(pipelinePlans) > 0 {
		c.Start()
	}
	return nil
}
