package services

import (
	"context"
	"github.com/merico-dev/lake/models/migrationscripts"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/migration"
	"github.com/merico-dev/lake/runner"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var cfg *viper.Viper
var db *gorm.DB
var cronManager *cron.Cron

func init() {
	var err error
	cfg = config.GetConfig()
	db, err = runner.NewGormDb(cfg, logger.Global.Nested("db"))
	location := cron.WithLocation(time.UTC)
	cronManager = cron.New(location)
	if err != nil {
		panic(err)
	}
	migration.Init(db)
	migrationscripts.RegisterAll()
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
	err = migration.Execute(context.Background())
	if err != nil {
		panic(err)
	}

	// call service init
	pipelineServiceInit()
}
