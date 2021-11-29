package tasks

import (
	"errors"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"gorm.io/gorm"
)

type DefaultJenkinsStorage struct {
	db *gorm.DB
}

func NewDefaultJenkinsStorage(db *gorm.DB) *DefaultJenkinsStorage {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get sql db", err)
		return nil
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return &DefaultJenkinsStorage{
		db,
	}
}

func (s *DefaultJenkinsStorage) SaveJob(job models.JenkinsJobProps) (context interface{}, err error) {
	var jenkinsJob = models.JenkinsJob{
		JenkinsJobProps: models.JenkinsJobProps{
			Name:  job.Name,
			Class: job.Class,
			Color: job.Color,
		},
	}
	var res = s.db.Save(&jenkinsJob)
	return jenkinsJob, res.Error
}

func (s *DefaultJenkinsStorage) SaveJobs(jobs []models.JenkinsJobProps) (context interface{}, err error) {
	var jenkinsJobs = make([]models.JenkinsJob, len(jobs))
	for index, job := range jobs {
		var jenkinsJob = models.JenkinsJob{
			JenkinsJobProps: models.JenkinsJobProps{
				Name:  job.Name,
				Class: job.Class,
				Color: job.Color,
			},
		}
		jenkinsJobs[index] = jenkinsJob
	}
	var res = s.db.Save(jenkinsJobs)
	return jenkinsJobs, res.Error
}

func (s *DefaultJenkinsStorage) SaveBuild(build models.JenkinsBuildProps, ctx interface{}) (context interface{}, err error) {
	var job, ok = ctx.(models.JenkinsJob)
	if !ok {
		return nil, errors.New("failed to get job id")
	}
	var jenkinsBuild = models.JenkinsBuild{
		JobID: job.ID,
		JenkinsBuildProps: models.JenkinsBuildProps{
			Duration:          build.Duration,
			DisplayName:       build.DisplayName,
			EstimatedDuration: build.EstimatedDuration,
			Number:            build.Number,
			Result:            build.Result,
			Timestamp:         build.Timestamp,
		},
	}
	var res = s.db.Save(&jenkinsBuild)
	return jenkinsBuild, res.Error
}

func (s *DefaultJenkinsStorage) SaveBuilds(builds []models.JenkinsBuildProps, ctx interface{}) (context interface{}, err error) {
	var job, ok = ctx.(models.JenkinsJob)
	if !ok {
		return nil, errors.New("failed to get job id")
	}
	var jenkinsBuilds = make([]models.JenkinsBuild, len(builds))
	for index, build := range builds {
		var jenkinsBuild = models.JenkinsBuild{
			JobID: job.ID,
			JenkinsBuildProps: models.JenkinsBuildProps{
				Duration:          build.Duration,
				DisplayName:       build.DisplayName,
				EstimatedDuration: build.EstimatedDuration,
				Number:            build.Number,
				Result:            build.Result,
				Timestamp:         build.Timestamp,
				StartTime:         time.Unix(build.Timestamp/1000, 0),
			},
		}
		jenkinsBuilds[index] = jenkinsBuild
	}
	var res = s.db.Save(&jenkinsBuilds)
	return jenkinsBuilds, res.Error
}
