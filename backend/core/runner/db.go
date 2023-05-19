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
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var _ config.ConfigReader = (*BaseDbConfigReader)(nil)

type BaseDbConfigReader struct {
	DbUrl          string `json:"DB_URL"`
	DbLoggingLevel string `json:"DB_LOGGING_LEVEL"`
	DbIdleConns    int    `json:"DB_IDLE_CONNS"`
	DbMaxConns     int    `json:"DB_MAX_CONNS"`
}

func (c *BaseDbConfigReader) Get(key string) interface{} {
	return nil
}

func (c *BaseDbConfigReader) GetBool(name string) bool {
	return false
}

func (c *BaseDbConfigReader) GetFloat64(key string) float64 {
	return 0
}

func (c *BaseDbConfigReader) GetInt(key string) int {
	switch key {
	case "DB_IDLE_CONNS":
		return c.DbIdleConns
	case "DB_MAX_CONNS":
		return c.DbMaxConns
	default:
		return 0
	}
}

func (c *BaseDbConfigReader) GetInt64(key string) int64 {
	return int64(c.GetInt(key))
}

func (c *BaseDbConfigReader) GetUint(key string) uint {
	return uint(c.GetInt(key))
}

func (c *BaseDbConfigReader) GetUint64(key string) uint64 {
	return uint64(c.GetInt(key))
}

func (c *BaseDbConfigReader) GetIntSlice(key string) []int {
	return []int{
		c.GetInt(key),
	}
}

func (c *BaseDbConfigReader) GetString(key string) string {
	switch key {
	case "DB_LOGGING_LEVEL":
		return c.DbLoggingLevel
	case "DB_URL":
		return c.DbUrl
	default:
		return ""
	}
}

func (c *BaseDbConfigReader) GetStringMap(key string) map[string]interface{} {
	return map[string]interface{}{
		key: c.Get(key),
	}
}

func (c *BaseDbConfigReader) GetStringMapString(key string) map[string]string {
	return map[string]string{
		key: c.GetString(key),
	}
}

func (c *BaseDbConfigReader) GetStringSlice(key string) []string {
	return []string{
		c.GetString(key),
	}
}

func (c *BaseDbConfigReader) GetTime(key string) time.Time {
	return time.Now()
}

func (c *BaseDbConfigReader) GetDuration(key string) time.Duration {
	return time.Duration(c.GetTime(key).Unix())
}

func (c *BaseDbConfigReader) IsSet(key string) bool {
	return false
}

func (c *BaseDbConfigReader) AllSettings() map[string]interface{} {
	return map[string]interface{}{}
}

// NewGormDb creates a new *gorm.DB and set it up properly
func NewGormDb(configReader config.ConfigReader, logger log.Logger) (*gorm.DB, errors.Error) {
	return NewGormDbEx(configReader, logger, &dal.SessionConfig{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
}

// NewGormDbEx acts like NewGormDb but accept extra sessionConfig
func NewGormDbEx(configReader config.ConfigReader, logger log.Logger, sessionConfig *dal.SessionConfig) (*gorm.DB, errors.Error) {
	dbLoggingLevel := gormLogger.Error
	switch strings.ToLower(configReader.GetString("DB_LOGGING_LEVEL")) {
	case "silent":
		dbLoggingLevel = gormLogger.Silent
	case "warn":
		dbLoggingLevel = gormLogger.Warn
	case "info":
		dbLoggingLevel = gormLogger.Info
	}

	idleConns := configReader.GetInt("DB_IDLE_CONNS")
	if idleConns <= 0 {
		idleConns = 10
	}
	dbMaxOpenConns := configReader.GetInt("DB_MAX_CONNS")
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
		PrepareStmt:            sessionConfig.PrepareStmt,
		SkipDefaultTransaction: sessionConfig.SkipDefaultTransaction,
	}
	dbUrl := configReader.GetString("DB_URL")
	if dbUrl == "" {
		return nil, errors.BadInput.New("DB_URL is required")
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, errors.Convert(err)
	}
	var db *gorm.DB
	switch strings.ToLower(u.Scheme) {
	case "mysql":
		dbUrl = fmt.Sprintf("%s@tcp(%s)%s?%s", getUserString(u), u.Host, u.Path, addLocal(u.Query()))
		db, err = gorm.Open(mysql.Open(dbUrl), dbConfig)
	case "postgresql", "postgres", "pg":
		db, err = gorm.Open(postgres.Open(dbUrl), dbConfig)
	default:
		return nil, errors.BadInput.New(fmt.Sprintf("invalid DB_URL:%s", dbUrl))
	}
	if err != nil {
		return nil, errors.Convert(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Convert(err)
	}
	sqlDB.SetMaxIdleConns(idleConns)
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, errors.Convert(err)
}

func getUserString(u *url.URL) string {
	userString := u.User.Username()
	password, ok := u.User.Password()
	if ok {
		userString = fmt.Sprintf("%s:%s", userString, password)
	}
	return userString
}

// addLocal adds loc=Local to the query string if it's not already there
func addLocal(query url.Values) string {
	if query.Get("loc") == "" {
		query.Set("loc", "Local")
	}
	return query.Encode()
}
