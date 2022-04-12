package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertApiMergeRequestsMeta = core.SubTaskMeta{
	Name:             "convertApiMergeRequests",
	EntryPoint:       ConvertApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Update domain layer PullRequest according to GitlabMergeRequest",
}

func ConvertApiMergeRequests(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	db := taskCtx.GetDb()

	domainMrIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	//Find all piplines associated with the current projectid
	cursor, err := db.Model(&models.GitlabMergeRequest{}).Where("project_id=?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequest{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabMr := inputRow.(*models.GitlabMergeRequest)

			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainMrIdGenerator.Generate(gitlabMr.GitlabId),
				},
				BaseRepoId:  domainRepoIdGenerator.Generate(gitlabMr.ProjectId),
				Status:      gitlabMr.State,
				Title:       gitlabMr.Title,
				Url:         gitlabMr.WebUrl,
				CreatedDate: gitlabMr.GitlabCreatedAt,
				MergedDate:  gitlabMr.MergedAt,
				ClosedDate:  gitlabMr.ClosedAt,
			}

			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
