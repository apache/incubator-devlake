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
	"sync"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models/migrationscripts"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

var cfg config.ConfigReader
var logger log.Logger
var db dal.Dal

var bpManager *services.BlueprintManager
var basicRes context.BasicRes
var migrator plugin.Migrator
var cronManager *cron.Cron
var vld *validator.Validate
var serviceStatus string

const (
	SERVICE_STATUS_INIT         = "initializing"
	SERVICE_STATUS_WAIT_CONFIRM = "waiting for migration confirmation"
	SERVICE_STATUS_MIGRATING    = "migrating"
	SERVICE_STATUS_READY        = "ready"
)
const failToCreateCronJob = "created cron job failed"

// InitResources creates resources needed by services module
func InitResources() {
	var err error

	// basic resources initialization
	serviceStatus = SERVICE_STATUS_INIT
	vld = validator.New()
	basicRes = runner.CreateAppBasicRes()
	cfg = basicRes.GetConfigReader()
	logger = basicRes.GetLogger()
	db = basicRes.GetDal()
	bpManager = services.NewBlueprintManager(db)
	// initialize db migrator
	migrator, err = runner.InitMigrator(basicRes)
	if err != nil {
		panic(err)
	}
	logger.Info("migration initialized")
	migrator.Register(migrationscripts.All(), "Framework")
}

// GetBasicRes returns the context.BasicRes instance used by services module
func GetBasicRes() context.BasicRes {
	return basicRes
}

// GetMigrator returns the core.Migrator instance used by services module
func GetMigrator() plugin.Migrator {
	return migrator
}

// Init the services module
// Should not be called concurrently
func Init() {
	InitResources()

	// lock the database to avoid multiple devlake instances from sharing the same one
	lockDatabase()

	// now, load the plugins
	errors.Must(runner.LoadPlugins(basicRes))

	// pull migration scripts from plugins to migrator
	for _, pluginInst := range plugin.AllPlugins() {
		if migratable, ok := pluginInst.(plugin.PluginMigration); ok {
			migrator.Register(migratable.MigrationScripts(), pluginInst.Name())
		}
	}

	// check if there are pending migration
	if migrator.HasPendingScripts() {
		if cfg.GetBool("FORCE_MIGRATION") {
			errors.Must(ExecuteMigration())
			logger.Info("db migration without confirmation")
		} else {
			serviceStatus = SERVICE_STATUS_WAIT_CONFIRM
			logger.Info("db migration confirmation needed")
		}
	} else {
		errors.Must(ExecuteMigration())
		logger.Info("no db migration needed")
	}
}

var statusLock sync.Mutex

// ExecuteMigration executes all pending migration scripts and initialize services module
// This might be called concurrently across multiple API requests
func ExecuteMigration() errors.Error {
	statusLock.Lock()
	defer statusLock.Unlock()
	if serviceStatus == SERVICE_STATUS_MIGRATING {
		return errors.BadInput.New("already migrating")
	}
	if serviceStatus == SERVICE_STATUS_READY {
		return nil
	}
	serviceStatus = SERVICE_STATUS_MIGRATING
	statusLock.Unlock() // unlock to allow other API requests to check the status
	// apply all pending migration scripts
	err := migrator.Execute()
	if err != nil {
		return err
	}

	// cronjob for blueprint triggering
	location := cron.WithLocation(time.UTC)
	cronManager = cron.New(location)

	// initialize pipeline server, mainly to start the pipeline consuming process
	pipelineServiceInit()
	statusLock.Lock()
	serviceStatus = SERVICE_STATUS_READY
	return nil
}

func CurrentStatus() string {
	return serviceStatus
}
