package tasks

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	models "github.com/apache/incubator-devlake/plugins/jenkins/models"
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
