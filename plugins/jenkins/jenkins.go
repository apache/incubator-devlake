package main

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jenkins/api"
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

func (j Jenkins) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var op = JenkinsOptions{
		Host:     config.V.GetString("JENKINS_ENDPOINT"),
		Username: config.V.GetString("JENKINS_USERNAME"),
		Password: config.V.GetString("JENKINS_PASSWORD"),
	}
	var err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("Failed to decode options: %v", err)
	}
	j.CleanData()
	var worker = tasks.NewJenkinsWorker(nil, tasks.NewDefaultJenkinsStorage(lakeModels.Db), op.Host, op.Username, op.Password)
	err = worker.SyncJobs(progress)
	if err != nil{
		logger.Error("Fail to sync jobs", err)
		return err
	}
	err = tasks.ConvertJobs()
	if err != nil{
		logger.Error("Fail to convert jobs", err)
		return err
	}
	err = tasks.ConvertBuilds()
	if err != nil{
		logger.Error("Fail to convert builds", err)
		return err
	}
	return nil
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jenkins"
}

func (plugin Jenkins) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
		"sources": {
			"GET":  api.ListSources,
			"POST": api.PostSource,
		},
		"sources/:sourceId": {
			"GET": api.GetSource,
			"PUT": api.PutSource,
		},
	}
}

var PluginEntry Jenkins //nolint
