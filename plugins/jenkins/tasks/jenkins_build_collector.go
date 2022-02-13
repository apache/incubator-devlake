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

func CollectBuilds(worker *JenkinsApiClient, scheduler *utils.WorkerScheduler) error {
	ctx := context.Background()
	cursor, err := lakeModels.Db.Model(&models.JenkinsJob{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	var builds struct {
		Builds []models.JenkinsBuildProps `json:"allBuilds"`
	}

	for cursor.Next() {
		jobCtx := &models.JenkinsJob{}
		err = lakeModels.Db.ScanRows(cursor, jobCtx)
		lastJenkinsBuild := &models.JenkinsBuild{}
		lakeModels.Db.Where("job_id = ?", jobCtx.ID).Order("timestamp DESC").First(lastJenkinsBuild)

		job, err := worker.jenkins.GetJob(ctx, jobCtx.Name)
		err = scheduler.Submit(func() error {
			logger.Debug("(collect build) Submit", job)
			_, err = job.Jenkins.Requester.GetJSON(ctx, job.Base, &builds,
				map[string]string{"tree": "allBuilds[number,timestamp,duration,estimatedDuration,displayName,result]"})

			if err != nil {
				return fmt.Errorf("failed to get jenkins job builds: %v", err)
			}
			// jenkins api is not supported to filter data with timestampe
			// so we need to filter it manually
			//timestampHalfYearAgo := time.Now().AddDate(0, -2, 0).Unix() * 1000

			var filteredData = make([]models.JenkinsBuildProps, 0)
			for _, item := range builds.Builds {
				if item.Timestamp > lastJenkinsBuild.Timestamp {
					filteredData = append(filteredData, item)
				}
			}

			//var filteredData = builds.Builds

			var jenkinsBuilds = make([]models.JenkinsBuild, len(filteredData))
			for index, build := range filteredData {
				var jenkinsBuild = models.JenkinsBuild{
					JobID: jobCtx.ID,
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
