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
	"time"

	"sync"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/impl/dalgorm"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/runner"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var cfg *viper.Viper
var log core.Logger
var db *gorm.DB
var basicRes core.BasicRes
var migrator core.Migrator
var cronManager *cron.Cron
var cronLocker sync.Mutex

const failToCreateCronJob = "created cron job failed"

// Init the services module
func Init() {
	var err error
	cfg = config.GetConfig()
	log = logger.Global
	db, err = runner.NewGormDb(cfg, log)
	if err != nil {
		panic(err)
	}

	// TODO: this is ugly, the lockDb / CreateAppBasicRes are coupled via global variables cfg/log
	// it is too much for this refactor, let's solve it later

	// lock the database to avoid multiple devlake instances from sharing the same one
	lockDb()

	// basic resources initialization
	basicRes = runner.CreateBasicRes(cfg, log, db)

	// initialize db migrator
	migrator, err = runner.InitMigrator(basicRes)
	if err != nil {
		panic(err)
	}
	log.Info("migration initialized")
	migrator.Register(migrationscripts.All(), "Framework")

	// now,
	// load plugins
	err = runner.LoadPlugins(basicRes)
	if err != nil {
		panic(err)
	}
	for pluginName, pluginInst := range core.AllPlugins() {
		if migratable, ok := pluginInst.(core.PluginMigration); ok {
			migrator.Register(migratable.MigrationScripts(), pluginName)
		}
	}
	forceMigration := cfg.GetBool("FORCE_MIGRATION")
	if !migrator.HasPendingScripts() || forceMigration {
		err = ExecuteMigration()
		if err != nil {
			panic(err)
		}
	}
	log.Info("Db migration confirmation needed")
}

// ExecuteMigration executes all pending migration scripts and initialize services module
func ExecuteMigration() errors.Error {
	// apply all pending migration scripts
	err := migrator.Execute()
	if err != nil {
		return err
	}

	// cronjob for blueprint triggering
	location := cron.WithLocation(time.UTC)
	cronManager = cron.New(location)
	if err != nil {
		panic(err)
	}
	// call service init
	pipelineServiceInit()
	return nil
}

// MigrationRequireConfirmation returns if there were migration scripts waiting to be executed
func MigrationRequireConfirmation() bool {
	return migrator.HasPendingScripts()
}

func lockDb() {
	// gorm doesn't support creating a PrepareStmt=false session from a PrepareStmt=true
	// but the lockDatabase needs PrepareStmt=false for table locking, we have to deal with it here
	lockingDb, err := runner.NewGormDbEx(cfg, logger.Global.Nested("migrator db"), &dal.SessionConfig{
		PrepareStmt:            false,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	err = lockDatabase(dalgorm.NewDalgorm(lockingDb))
	if err != nil {
		panic(err)
	}
}
