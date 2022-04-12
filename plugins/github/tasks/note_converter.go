package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertNotesMeta = core.SubTaskMeta{
	Name:             "convertNotes",
	EntryPoint:       ConvertNotes,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_comments into  domain layer table notes",
}

func ConvertNotes(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	noteIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequestComment{})
	userIdGen := didgen.NewDomainIdGenerator(&models.GithubUser{})
	prIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})

	cursor, err := db.Model(&models.GithubPullRequestComment{}).
		Joins(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_comments.pull_request_id`).
		Where("_tool_github_pull_requests.repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubPullRequestComment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			prComment := inputRow.(*models.GithubPullRequestComment)
			domainNote := &code.Note{
				DomainEntity: domainlayer.DomainEntity{
					Id: noteIdGen.Generate(prComment.GithubId),
				},
				PrId:        prIdGen.Generate(prComment.PullRequestId),
				Author:      userIdGen.Generate(prComment.AuthorUserId),
				Body:        prComment.Body,
				CreatedDate: prComment.GithubCreatedAt,
			}
			return []interface{}{
				domainNote,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
