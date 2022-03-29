package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
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

	noteIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequestComment{})

	cursor, err := db.Model(&githubModels.GithubPullRequestComment{}).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_comments.pull_request_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestComment{}),
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
			prComment := inputRow.(*githubModels.GithubPullRequestComment)
			domainNote := &code.Note{
				DomainEntity: domainlayer.DomainEntity{
					Id: noteIdGen.Generate(prComment.GithubId),
				},
				PrId:        uint64(prComment.PullRequestId),
				Author:      prComment.AuthorUsername,
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
