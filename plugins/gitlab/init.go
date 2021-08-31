package main

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/gitlab/models"
)

func (plugin Gitlab) Init() {
	logger.Info("INFO >>> init go plugin", true)
	err := lakeModels.Db.AutoMigrate(&models.GitlabCommit{}, &models.GitlabProject{})
	if err != nil {
		logger.Error("Error migrating gitlab: ", err)
		panic(err)
	}
}
