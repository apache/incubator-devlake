package services

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/runner"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var cfg *viper.Viper
var db *gorm.DB

func init() {
	var err error
	cfg = config.GetConfig()
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
	temporalUrl := cfg.GetString("TEMPORAL_URL")
	if temporalUrl == "" {
		err = runner.MigrateDb(db)
		if err != nil {
			panic(err)
		}
	} else {

	}

	// call service init
	pipelineServiceInit()
	taskServiceInit()
}
