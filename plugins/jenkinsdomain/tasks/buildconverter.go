package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	domainlayerBase "github.com/merico-dev/lake/plugins/domainlayer/models/base"
	"github.com/merico-dev/lake/plugins/domainlayer/models/devops"
	"github.com/merico-dev/lake/plugins/domainlayer/okgen"
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

	jobOriginkeyGenerator := okgen.NewOriginKeyGenerator(&jenkinsModels.JenkinsJob{})
	buildOriginkeyGenerator := okgen.NewOriginKeyGenerator(jenkinsBuild)

	println("convert builds")
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jenkinsBuild)
		if err != nil {
			return err
		}
		println("jenkins build", jenkinsBuild)
		build := &devops.Build{
			DomainEntity: domainlayerBase.DomainEntity{
				OriginKey: buildOriginkeyGenerator.Generate(jenkinsBuild.ID),
			},
			JobOriginKey: jobOriginkeyGenerator.Generate(jenkinsBuild.JobID),
			Name:         jenkinsBuild.DisplayName,
			DurationSec:  uint64(jenkinsBuild.Duration),
			Status:       jenkinsBuild.Result,
			StartedDate:  jenkinsBuild.StartTime,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(build).Error
		if err != nil {
			return err
		}
	}
	return nil
}
