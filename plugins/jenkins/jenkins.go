package main

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"github.com/merico-dev/lake/plugins/jenkins/tasks"
	"github.com/mitchellh/mapstructure"
)

type JenkinsOptions struct {
	Host     string
	Username string
	Password string
}

type Jenkins struct{}

func (j Jenkins) Init() {
	var err = lakeModels.Db.AutoMigrate(&models.JenkinsJob{}, &models.JenkinsBuild{})
	if err != nil {
		logger.Error("Failed to auto migrate jenkins models", err)
	}
}

func (j Jenkins) Description() string {
	return "Jenkins plugin"
}

func (j Jenkins) Execute(options map[string]interface{}, progress chan<- float32) {
	var op JenkinsOptions
	var err = mapstructure.Decode(options, &op)
	if err != nil {
		logger.Error("Failed to decode options", err)
		return
	}
	var worker = tasks.NewJenkinsWorker(nil, tasks.NewDeafultJenkinsStorage(lakeModels.Db), op.Host, op.Username, op.Password)
	worker.SyncJobs()
}

var PluginEntry Jenkins
