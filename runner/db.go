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

package runner

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewGormDb FIXME ...
func NewGormDb(config *viper.Viper, logger core.Logger) (*gorm.DB, error) {
	dbLoggingLevel := gormLogger.Error
	switch strings.ToLower(config.GetString("DB_LOGGING_LEVEL")) {
	case "silent":
		dbLoggingLevel = gormLogger.Silent
	case "warn":
		dbLoggingLevel = gormLogger.Warn
	case "info":
		dbLoggingLevel = gormLogger.Info
	}

	idleConns := config.GetInt("DB_IDLE_CONNS")
	if idleConns <= 0 {
		idleConns = 10
	}
	dbMaxOpenConns := config.GetInt("DB_MAX_CONNS")
	if dbMaxOpenConns <= 0 {
		dbMaxOpenConns = 100
	}

	dbConfig := &gorm.Config{
		Logger: gormLogger.New(
			logger,
			gormLogger.Config{
				SlowThreshold:             time.Second,    // Slow SQL threshold
				LogLevel:                  dbLoggingLevel, // Log level
				IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,           // Disable color
			},
		),
		// most of our operation are in batch, this can improve performance
		PrepareStmt: true,
	}
	dbUrl := config.GetString("DB_URL")
	if dbUrl == "" {
		return nil, errors.BadInput.New("DB_URL is required", errors.AsUserMessage())
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	var db *gorm.DB
	switch strings.ToLower(u.Scheme) {
	case "mysql":
		dbUrl = fmt.Sprintf("%s@tcp(%s)%s?%s", u.User.String(), u.Host, u.Path, u.RawQuery)
		db, err = gorm.Open(mysql.Open(dbUrl), dbConfig)
	case "postgresql", "postgres", "pg":
		db, err = gorm.Open(postgres.Open(dbUrl), dbConfig)
	default:
		return nil, errors.BadInput.New(fmt.Sprintf("invalid DB_URL:%s", dbUrl), errors.AsUserMessage())
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(idleConns)
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, err
}
