package main

import (
	"github.com/merico-dev/lake/config"
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

func (j Jenkins) CleanData() {
	var err = lakeModels.Db.Exec("truncate table jenkins_jobs").Error
	if err != nil {
		logger.Error("Failed to truncate jenkins models", err)
	}
	err = lakeModels.Db.Exec("truncate table jenkins_builds").Error
	if err != nil {
		logger.Error("Failed to truncate jenkins models", err)
	}
}

func (j Jenkins) Execute(options map[string]interface{}, taskId uint64, progress chan<- float32) {
	var op = JenkinsOptions{
		Host:     config.V.GetString("JENKINS_ENDPOINT"),
		Username: config.V.GetString("JENKINS_USERNAME"),
		Password: config.V.GetString("JENKINS_PASSWORD"),
	}
	logger.Info("Jenkins config", op)
	var err = mapstructure.Decode(options, &op)
	if err != nil {
		logger.Error("Failed to decode options", err)
		return
	}
	j.CleanData()
	var worker = tasks.NewJenkinsWorker(nil, tasks.NewDeafultJenkinsStorage(lakeModels.Db), op.Host, op.Username, op.Password)
	worker.SyncJobs(progress)
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jenkins"
}

var PluginEntry Jenkins //nolint
