package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	jenkinsModels "github.com/merico-dev/lake/plugins/jenkins/models"
	"gorm.io/gorm/clause"
)

func ConvertBuilds() error {
	jenkinsBuild := &jenkinsModels.JenkinsBuild{}

	cursor, err := lakeModels.Db.Model(jenkinsBuild).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&jenkinsModels.JenkinsJob{})
	buildIdGen := didgen.NewDomainIdGenerator(jenkinsBuild)

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jenkinsBuild)
		if err != nil {
			return err
		}
		build := &devops.Build{
			DomainEntity: domainlayer.DomainEntity{
				Id: buildIdGen.Generate(jenkinsBuild.ID),
			},
			JobId:       jobIdGen.Generate(jenkinsBuild.JobID),
			Name:        jenkinsBuild.DisplayName,
			DurationSec: uint64(jenkinsBuild.Duration),
			Status:      jenkinsBuild.Result,
			StartedDate: jenkinsBuild.StartTime,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(build).Error
		if err != nil {
			return err
		}
	}
	return nil
}
