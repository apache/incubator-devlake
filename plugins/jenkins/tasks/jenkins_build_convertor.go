package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	models "github.com/merico-dev/lake/plugins/jenkins/models"
	"reflect"
)

var ConvertBuildsMeta = core.SubTaskMeta{
	Name:             "convertBuilds",
	EntryPoint:       ConvertBuilds,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_builds into  domain layer table builds",
}

func ConvertBuilds(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()

	jenkinsBuild := &models.JenkinsBuild{}

	cursor, err := db.Model(jenkinsBuild).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})
	buildIdGen := didgen.NewDomainIdGenerator(&models.JenkinsBuild{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsBuild{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jenkinsBuild := inputRow.(*models.JenkinsBuild)
			build := &devops.Build{
				DomainEntity: domainlayer.DomainEntity{
					Id: buildIdGen.Generate(jenkinsBuild.JobName, jenkinsBuild.Number),
				},
				JobId:       jobIdGen.Generate(jenkinsBuild.JobName),
				Name:        jenkinsBuild.DisplayName,
				DurationSec: uint64(jenkinsBuild.Duration / 1000),
				Status:      jenkinsBuild.Result,
				StartedDate: jenkinsBuild.StartTime,
				CommitSha:   jenkinsBuild.CommitSha,
			}
			return []interface{}{
				build,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
