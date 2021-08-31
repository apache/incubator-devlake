package tasks

import (
	"github.com/bndr/gojenkins"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"gorm.io/gorm"
)

type JenkinsStorage interface {
	SaveJob(job *gojenkins.Job) error
	SaveJobs(jobs []*gojenkins.Job) error
	SaveBuild(build *gojenkins.Build) error
	SaveBuilds(build []*gojenkins.Build) error
}

type DefaultJenkinsStorage struct {
	db *gorm.DB
}

func NewDeafultJenkinsStorage(db *gorm.DB) *DefaultJenkinsStorage {
	return &DefaultJenkinsStorage{
		db,
	}
}

func (s *DefaultJenkinsStorage) SaveJob(job *gojenkins.Job) error {
	var jenkinsJob = models.JenkinsJob{
		Name:  job.Raw.Name,
		Class: job.Raw.Class,
		Color: job.Raw.Color,
	}
	return s.db.Save(&jenkinsJob).Error
}

func (s *DefaultJenkinsStorage) SaveJobs(jobs []*gojenkins.Job) error {
	var jenkinsJobs = make([]models.JenkinsJob, len(jobs))
	for index, item := range jobs {
		var jenkinsJob = models.JenkinsJob{
			Name:  item.Raw.Name,
			Class: item.Raw.Class,
			Color: item.Raw.Color,
		}
		jenkinsJobs[index] = jenkinsJob
	}
	return s.db.Save(jenkinsJobs).Error
}

func (s *DefaultJenkinsStorage) SaveBuild(build *gojenkins.Build) error {
	var jenkinsBuild = models.JenkinsBuild{
		Duration:          build.Raw.Duration,
		DisplayName:       build.Raw.DisplayName,
		EstimatedDuration: build.Raw.EstimatedDuration,
		Number:            build.Raw.Number,
		Result:            build.Raw.Result,
		Timestamp:         build.Raw.Timestamp,
	}
	return s.db.Save(&jenkinsBuild).Error
}

func (s *DefaultJenkinsStorage) SaveBuilds(build []*gojenkins.Build) error {
	var jenkinsBuilds = make([]models.JenkinsBuild, len(build))
	for index, item := range build {
		var jenkinsBuild = models.JenkinsBuild{
			Duration:          item.Raw.Duration,
			DisplayName:       item.Raw.DisplayName,
			EstimatedDuration: item.Raw.EstimatedDuration,
			Number:            item.Raw.Number,
			Result:            item.Raw.Result,
			Timestamp:         item.Raw.Timestamp,
		}
		jenkinsBuilds[index] = jenkinsBuild
	}
	return s.db.Save(jenkinsBuilds).Error
}
