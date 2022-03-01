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

func ConvertBuilds(ctx context.Context) error {
	err := lakeModels.Db.Where("id like 'jenkins:JenkinsBuild:%'").Delete(&devops.Build{}).Error
	if err != nil {
		return err
	}

	jenkinsBuild := &jenkinsModels.JenkinsBuild{}

	cursor, err := lakeModels.Db.Model(jenkinsBuild).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&jenkinsModels.JenkinsJob{})
	buildIdGen := didgen.NewDomainIdGenerator(&jenkinsModels.JenkinsBuild{})

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jenkinsBuild)
		if err != nil {
			return err
		}
		build := &devops.Build{
			DomainEntity: domainlayer.DomainEntity{
				Id: buildIdGen.Generate(jenkinsBuild.JobName, jenkinsBuild.Number),
			},
			JobId:       jobIdGen.Generate(jenkinsBuild.JobName),
			Name:        jenkinsBuild.DisplayName,
			DurationSec: uint64(jenkinsBuild.Duration),
			Status:      jenkinsBuild.Result,
			StartedDate: jenkinsBuild.StartTime,
			CommitSha:   jenkinsBuild.CommitSha,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(build).Error
		if err != nil {
			return err
		}
	}
	return nil
}
