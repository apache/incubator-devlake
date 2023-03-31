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

package helper

import (
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/impls/logruslog"
)

// InitDB Bootstraps the database by getting rid of all the tables
func InitDB(dbUrl string) {
	logger := logruslog.Global.Nested("test-init")
	logger.Info("Initializing database")
	cfg := config.GetConfig()
	cfg.Set("DB_URL", dbUrl)
	db, err := runner.NewGormDb(cfg, logger)
	if err != nil {
		panic(err)
	}
	migrator := db.Migrator()
	tables, err := errors.Convert01(migrator.GetTables())
	if err != nil {
		panic(err)
	}
	logger.Info("Dropping %d existing tables", len(tables))
	var tablesRaw []any
	for _, table := range tables {
		tablesRaw = append(tablesRaw, table)
	}
	err = errors.Convert(migrator.DropTable(tablesRaw...))
	if err != nil {
		panic(err)
	}
}
