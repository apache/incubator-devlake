package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	jenkinsModels "github.com/merico-dev/lake/plugins/jenkins/models"
	"gorm.io/gorm/clause"
)

func ConvertJobs(ctx context.Context) error {
	jenkinsJob := &jenkinsModels.JenkinsJob{}

	jobIdGen := didgen.NewDomainIdGenerator(jenkinsJob)
	err := lakeModels.Db.
		Delete(&devops.Job{}, "`name` not in (select `name` from jenkins_jobs)").Error
	if err != nil {
		return err
	}

	cursor, err := lakeModels.Db.Model(jenkinsJob).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jenkinsJob)
		if err != nil {
			return err
		}
		job := &devops.Job{
			DomainEntity: domainlayer.DomainEntity{
				Id: jobIdGen.Generate(jenkinsJob.Name),
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
