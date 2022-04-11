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

var ConvertMergeRequestCommentMeta = core.SubTaskMeta{
	Name:             "convertMergeRequestComment",
	EntryPoint:       ConvertMergeRequestComment,
	EnabledByDefault: true,
	Description:      "Update domain layer Comment according to GitlabMergeRequestComment",
}

func ConvertMergeRequestComment(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDb()

	cursor, err := db.Model(&models.GitlabMergeRequestComment{}).
		Joins("left join gitlab_merge_requests on gitlab_merge_requests.gitlab_id = gitlab_merge_request_comments.merge_request_id").
		Where("gitlab_merge_requests.project_id = ?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorComment := didgen.NewDomainIdGenerator(&models.GitlabMergeRequestComment{})
	prIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	userIdGen := didgen.NewDomainIdGenerator(&models.GitlabUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequestComment{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabComments := inputRow.(*models.GitlabMergeRequestComment)

			domainComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainIdGeneratorComment.Generate(gitlabComments.GitlabId),
				},
				PullRequestId: prIdGen.Generate(gitlabComments.MergeRequestId),
				Body:          gitlabComments.Body,
				UserId:        userIdGen.Generate(gitlabComments.AuthorUsername),
				CreatedDate:   gitlabComments.GitlabCreatedAt,
			}
			return []interface{}{
				domainComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
