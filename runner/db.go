package runner

import (
	"fmt"
	"net/url"
	"os/user"
	"time"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewGormDb(config *viper.Viper, logger core.Logger) (*gorm.DB, error) {
	dbUrl := config.GetString("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "mysql" {
		dbUrl = fmt.Sprintf(("%s@tcp(%s)%s?%s"), u.User.String(), u.Host, u.Path, u.RawQuery)
	}

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

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{
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
	})
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

func MigrateDb(db *gorm.DB) error {
	// domain layer entity
	return db.AutoMigrate(
		&user.User{},
		&code.Repo{},
		&code.Commit{},
		&code.CommitParent{},
		&code.PullRequest{},
		&code.PullRequestCommit{},
		&code.PullRequestLabel{},
		&code.Note{},
		&code.RepoCommit{},
		&code.Ref{},
		&code.RefsCommitsDiff{},
		&code.CommitFile{},
		&code.PullRequestComment{},
		&ticket.Board{},
		&ticket.Issue{},
		&ticket.IssueLabel{},
		&ticket.BoardIssue{},
		&ticket.BoardSprint{},
		&ticket.Changelog{},
		&ticket.Sprint{},
		&ticket.SprintIssue{},
		&ticket.IssueStatusHistory{},
		&ticket.IssueSprintsHistory{},
		&ticket.IssueAssigneeHistory{},
		&devops.Job{},
		&devops.Build{},
		&ticket.IssueWorklog{},
		&ticket.IssueComment{},
		&crossdomain.BoardRepo{},
		&crossdomain.PullRequestIssue{},
		&crossdomain.IssueCommit{},
		&crossdomain.IssueRepoCommit{},
		&crossdomain.RefsIssuesDiffs{},
		&code.RefsPrCherrypick{},
	)
}
