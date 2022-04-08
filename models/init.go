package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/models/domainlayer/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {
	V := config.GetConfig()
	connectionString := V.GetString("DB_URL")
	if V.Get("TEST") == "true" {
		connectionString = "merico:merico@tcp(localhost:3306)/lake_test"
	}
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	Db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("ERROR: >>> Mysql failed to connect")
		panic(err)
	}
	sqlDB, err := Db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 24)
	migrateDB()
}

func migrateDB() {
	err := Db.AutoMigrate(
		&Task{},
		&Notification{},
		&Pipeline{},
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
		&ticket.Worklog{},
		&crossdomain.BoardRepo{},
		&crossdomain.PullRequestIssue{},
		&crossdomain.IssueCommit{},
		&crossdomain.RefsIssuesDiffs{},
	)
	if err != nil {
		panic(err)
	}
	// TODO: create customer migration here
}
