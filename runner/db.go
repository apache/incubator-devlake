package runner

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewGormDb(config *viper.Viper, logger core.Logger) (*gorm.DB, error) {
	dbLoggingLevel := gormLogger.Error
	switch config.GetString("DB_LOGGING_LEVEL") {
	case "Silent":
		dbLoggingLevel = gormLogger.Silent
	case "Warn":
		dbLoggingLevel = gormLogger.Warn
	case "Info":
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
		return nil, fmt.Errorf("DB_URL is required")
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	var db *gorm.DB
	switch strings.ToLower(u.Scheme) {
	case "mysql":
		dbUrl = fmt.Sprintf(("%s@tcp(%s)%s?%s"), u.User.String(), u.Host, u.Path, u.RawQuery)
		db, err = gorm.Open(mysql.Open(dbUrl), dbConfig)
	case "postgresql", "postgres", "pg":
		db, err = gorm.Open(postgres.Open(dbUrl), dbConfig)
	default:
		return nil, fmt.Errorf("invalid DB_URL:%s", dbUrl)
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
