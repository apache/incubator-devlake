package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
	"time"
)

func CollectBuilds(apiClient *JenkinsApiClient, scheduler *utils.WorkerScheduler, ctx context.Context) error {
	err := lakeModels.Db.Delete(&models.JenkinsBuild{}, "job_name not in (select `name` from jenkins_jobs)").Error
	if err != nil {
		return err
	}

	cursor, err := lakeModels.Db.Model(&models.JenkinsJob{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		jobCtx := &models.JenkinsJob{}
		var builds struct {
			Builds []models.JenkinsBuildProps `json:"allBuilds"`
		}
		err = lakeModels.Db.ScanRows(cursor, jobCtx)
		if err != nil {
			return err
		}
		lastJenkinsBuild := &models.JenkinsBuild{}
		err = lakeModels.Db.Where("job_name = ?", jobCtx.Name).Order("timestamp DESC").Limit(1).Find(lastJenkinsBuild).Error
		if err != nil {
			return err
		}

		err = scheduler.Submit(func() error {
			job, err := apiClient.jenkins.GetJob(ctx, jobCtx.Name)
			if err != nil {
				return err
			}
			logger.Debug("(collect build) Submit", job)
			_, err = job.Jenkins.Requester.GetJSON(ctx, job.Base, &builds,
				map[string]string{"tree": "allBuilds[number,timestamp,duration,estimatedDuration,displayName,result]"})

			if err != nil {
				return fmt.Errorf("failed to get jenkins job builds: %v", err)
			}

			var filteredData = make([]models.JenkinsBuildProps, 0)
			for _, v := range builds.Builds {
				if v.Timestamp > lastJenkinsBuild.Timestamp {
					build, err := job.GetBuild(ctx, v.Number)
					if err != nil {
						return fmt.Errorf("failed to get jenkins build: %v, %s:%d", err, job.Raw.Name, v.Number)
					}
					logger.Debug("(collect build commit sha)", build.GetBuildNumber())

					v.CommitSha = build.GetRevision()
					filteredData = append(filteredData, v)
				}
			}

			var jenkinsBuilds = make([]models.JenkinsBuild, len(filteredData))
			for i, v := range filteredData {
				var jenkinsBuild = models.JenkinsBuild{
					JobName:           jobCtx.Name,
					Duration:          v.Duration,
					DisplayName:       v.DisplayName,
					EstimatedDuration: v.EstimatedDuration,
					Number:            v.Number,
					Result:            v.Result,
					Timestamp:         v.Timestamp,
					StartTime:         time.Unix(v.Timestamp/1000, 0),
					CommitSha:         v.CommitSha,
				}
				jenkinsBuilds[i] = jenkinsBuild
			}
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).CreateInBatches(&jenkinsBuilds, len(jenkinsBuilds)).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	scheduler.WaitUntilFinish()

	return nil
}
