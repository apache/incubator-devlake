package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	jenkinsModels "github.com/merico-dev/lake/plugins/jenkins/models"
	"gorm.io/gorm/clause"
)

func ConvertJobs() error {
	jenkinsJob := &jenkinsModels.JenkinsJob{}

	cursor, err := lakeModels.Db.Model(jenkinsJob).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobOriginkeyGenerator := okgen.NewOriginKeyGenerator(jenkinsJob)

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jenkinsJob)
		if err != nil {
			return err
		}
		job := &devops.Job{
			DomainEntity: domainlayer.DomainEntity{
				OriginKey: jobOriginkeyGenerator.Generate(jenkinsJob.ID),
			},
			Name: jenkinsJob.Name,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(job).Error
		if err != nil {
			return err
		}
	}
	return nil
}
