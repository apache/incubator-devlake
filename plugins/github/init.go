package main

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
)

func (plugin Github) Init() {
	logger.Info("INFO >>> init GitHub plugin", true)
	err := lakeModels.Db.AutoMigrate(
		&models.GithubRepository{},
		&models.GithubCommit{},
		&models.GithubPullRequest{},
		&models.GithubReviewer{},
		&models.GithubComment{},
	)
	if err != nil {
		logger.Error("Error migrating github: ", err)
		panic(err)
	}
}
