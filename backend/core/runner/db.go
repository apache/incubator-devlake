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
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	tlsMysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

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
		// PrepareStmt:            sessionConfig.PrepareStmt,
		SkipDefaultTransaction: sessionConfig.SkipDefaultTransaction,
	}
	dbUrl := configReader.GetString("DB_URL")
	if dbUrl == "" {
		return nil, errors.BadInput.New("DB_URL is required, please set it in environment variable or .env file")
	}
	db, err := getDbConnection(dbUrl, dbConfig)
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

// sanitizeQuery add default value to query and remove ca-cert from query
func sanitizeQuery(query url.Values) string {
	if query.Get("loc") == "" {
		query.Set("loc", "Local")
	}
	if query.Get("ca-cert") != "" {
		query.Del("ca-cert")
	}
	return query.Encode()
}

func getDbConnection(dbUrl string, conf *gorm.Config) (*gorm.DB, error) {
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(u.Scheme) {
	case "mysql":
		dbUrl = fmt.Sprintf("%s@tcp(%s)%s?%s", getUserString(u), u.Host, u.Path, sanitizeQuery(u.Query()))
		if u.Query().Get("ca-cert") != "" {
			rootCertPool := x509.NewCertPool()
			pem, err := os.ReadFile(u.Query().Get("ca-cert"))
			if err != nil {
				return nil, err
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				return nil, err
			}
			err = tlsMysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
			if err != nil {
				return nil, err
			}
			dbUrl = fmt.Sprintf("%s&tls=custom", dbUrl)
			db, err := sql.Open("mysql", dbUrl)
			if err != nil {
				return nil, err
			}
			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn: db,
			}), &gorm.Config{})

			return gormDB, err
		}
		return gorm.Open(mysql.Open(dbUrl), conf)
	case "postgresql", "postgres", "pg":
		return gorm.Open(postgres.Open(dbUrl), conf)
	default:
		return nil, fmt.Errorf("invalid DB_URL:%s", dbUrl)
	}
}

func CheckDbConnection(dbUrl string, d time.Duration) errors.Error {
	db, err := getDbConnection(dbUrl, &gorm.Config{})
	if err != nil {
		return errors.Convert(err)
	}
	ctx := context.Background()
	if d > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d)
		defer cancel()
	}
	return errors.Convert(db.WithContext(ctx).Exec("SELECT 1").Error)
}
