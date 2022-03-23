package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"reflect"
)

var ConvertJobsMeta = core.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_jobs into  domain layer table jobs",
}

func ConvertJobs(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()

	jenkinsJob := &models.JenkinsJob{}

	jobIdGen := didgen.NewDomainIdGenerator(jenkinsJob)

	cursor, err := db.Model(jenkinsJob).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsJob{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jenkinsJob := inputRow.(*models.JenkinsJob)
			job := &devops.Job{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobIdGen.Generate(jenkinsJob.Name),
				},
				Name: jenkinsJob.Name,
			}
			return []interface{}{
				job,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
