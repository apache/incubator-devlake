package services

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/runner"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var cfg *viper.Viper
var db *gorm.DB
var cronManager *cron.Cron

func init() {
	var err error
	cfg = config.GetConfig()
	db, err = runner.NewGormDb(cfg, logger.Global.Nested("db"))
	cronManager = cron.New()
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
		&models.Blueprint{},
	)
	if err != nil {
		panic(err)
	}

	// migrate data tables if run in standalone mode
	err = runner.MigrateDb(db)
	if err != nil {
		panic(err)
	}

	// call service init
	pipelineServiceInit()
}
