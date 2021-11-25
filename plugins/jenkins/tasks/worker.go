package tasks

import (
	"context"
	"fmt"
	"net/http"

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

func (worker *JenkinsWorker) SyncJobs(progress chan<- float32) error {
	var ctx = context.Background()
	// get all jobs
	var jobs, err = worker.jenkins.GetAllJobs(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get jobs from jenkins: %v", err)
	}
	for index, job := range jobs {
		err = worker.syncJob(ctx, job)
		if err != nil {
			return err
		}
		progress <- float32((index + 1)) / float32(len(jobs))
	}
	return nil
}

func (worker *JenkinsWorker) syncJob(ctx context.Context, job *gojenkins.Job) error {
	logger.Info("syncJob", job.Raw.Name)
	jobCtx, err := worker.storage.SaveJob(models.JenkinsJobProps{
		Name:  job.Raw.Name,
		Class: job.Raw.Class,
		Color: job.Raw.Color,
	})
	if err != nil {
		return fmt.Errorf("failed to save job: %v", err)
	}
	var builds struct {
		Builds []models.JenkinsBuildProps `json:"allBuilds"`
	}
	_, err = job.Jenkins.Requester.GetJSON(ctx, job.Base, &builds, map[string]string{"tree": "allBuilds[number,timestamp,duration,estimatedDuration,displayName,result]"})
	if err != nil {
		return fmt.Errorf("failed to get jenkins job builds: %v", err)
	}
	// jenkins api is not supported to filter data with timestampe
	// so we need to filter it manually
	// timestampHalfYearAgo := time.Now().AddDate(0, -6, 0).Unix() * 1000
	// var filteredData = make([]models.JenkinsBuildProps, 0)
	// for _, item := range builds.Builds {
	// 	if item.Timestamp >= timestampHalfYearAgo {
	// 		filteredData = append(filteredData, item)
	// 	}
	// }
	var filteredData = builds.Builds
	if len(filteredData) == 0 {
		return fmt.Errorf("job has no build exists: %v", job.Raw.Name)
	}
	_, err = worker.storage.SaveBuilds(filteredData, jobCtx)
	if err != nil {
		return fmt.Errorf("failed to save builds: %v", err)
	}
	return nil
}
