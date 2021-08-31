package tasks

import (
	"context"
	"net/http"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jenkins/models"
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
	for _, job := range jobs {
		worker.syncJob(ctx, job)
	}
}

func (worker *JenkinsWorker) syncJob(ctx context.Context, job *gojenkins.Job) {
	jobCtx, err := worker.storage.SaveJob(models.JenkinsJobProps{
		Name:  job.Raw.Name,
		Class: job.Raw.Class,
		Color: job.Raw.Color,
	})
	if err != nil {
		logger.Error("failed to save job", err)
		return
	}
	var builds struct {
		Builds []models.JenkinsBuildProps `json:"allBuilds"`
	}
	_, err = job.Jenkins.Requester.GetJSON(ctx, job.Base, &builds, map[string]string{"tree": "allBuilds[number,timestamp,duration,estimatedDuration,displayName,result]"})
	if err != nil {
		logger.Error("failed to get jenkins job builds", err)
		return
	}
	// jenkins api is not supported to filter data with timestampe
	// so we need to filter it manually
	timestampHalfYearAgo := time.Now().AddDate(0, -6, 0).Unix() * 1000
	var filteredData = make([]models.JenkinsBuildProps, 0)
	for _, item := range builds.Builds {
		if item.Timestamp >= timestampHalfYearAgo {
			filteredData = append(filteredData, item)
		}
	}
	if len(filteredData) == 0 {
		logger.Info("job has no build exists in past half year", job.Raw.Name)
		return
	}
	worker.storage.SaveBuilds(filteredData, jobCtx)
}
