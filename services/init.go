package services

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/runner"
	"gorm.io/gorm"
	"strings"
)

var db *gorm.DB

func init() {
	var err error
	cfg := config.GetConfig()
	db, err = runner.NewGormDb(cfg, logger.Global.Nested("db"))

	if err != nil {
		panic(err)
	}

	// load plugins
	err = runner.LoadPlugins(
		cfg.GetString("PLUGIN_DIR"),
		cfg,
		logger.Global.Nested("plugin"),
		db,
	)
	if err != nil {
		panic(err)
	}

	// migrate framework tables
	err = db.AutoMigrate(
		&models.Task{},
		&models.Notification{},
		&models.Pipeline{},
		&models.PipelinePlan{},
	)
	if err != nil {
		panic(err)
	}

	// migrate data tables if run in standalone mode
	if cfg.GetBool("STAND_ALONE") {
		err = runner.MigrateDb(db)
		if err != nil {
			panic(err)
		}
	}

	// set all unfinished tasks to failed
	db.Model(&models.Task{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)

	var notificationEndpoint = cfg.GetString("NOTIFICATION_ENDPOINT")
	var notificationSecret = cfg.GetString("NOTIFICATION_SECRET")
	if strings.TrimSpace(notificationEndpoint) != "" {
		notificationService = NewNotificationService(notificationEndpoint, notificationSecret)
	}
	db.Model(&models.Pipeline{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)

	err = ReloadPipelinePlans()
	if err != nil {
		panic(err)
	}
}
