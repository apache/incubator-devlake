package tasks

import (
	"context"
	"net/http"

	"github.com/bndr/gojenkins"
	"github.com/merico-dev/lake/logger"
)

type JenkinsWorker struct {
	jenkins *gojenkins.Jenkins
	storage JenkinsStorage
}

func NewJenkinsWorker(client *http.Client, storage JenkinsStorage, base string, auth ...interface{}) *JenkinsWorker {
	return &JenkinsWorker{
		storage: storage,
		jenkins: gojenkins.CreateJenkins(client, base, auth...),
	}
}

func (worker *JenkinsWorker) SyncJobs() {
	var ctx = context.Background()
	// get all jobs
	var jobs, err = worker.jenkins.GetAllJobs(ctx)
	if err != nil {
		logger.Error("Failed to get jobs from jenkins", err)
		return
	}
	err = worker.storage.SaveJobs(jobs)
	if err != nil {
		logger.Error("Failed to save jobs", err)
		return
	}
	// TODO: get builds
}
