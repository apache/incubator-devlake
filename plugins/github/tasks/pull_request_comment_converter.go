package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

var ConvertPullRequestCommentsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestComments",
	EntryPoint:       ConvertPullRequestComments,
	EnabledByDefault: true,
	Description:      "ConvertPullRequestComments data from Github api",
}

func ConvertPullRequestComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&githubModels.GithubPullRequestComment{}).
		Joins("left join github_pull_requests "+
			"on github_pull_requests.github_id = github_pull_request_comments.pull_request_id").
		Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})
	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

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
			githubPullRequestComment := inputRow.(*githubModels.GithubPullRequestComment)
			domainPrComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: prIdGen.Generate(githubPullRequestComment.GithubId),
				},
				PullRequestId: prIdGen.Generate(githubPullRequestComment.PullRequestId),
				Body:          githubPullRequestComment.Body,
				UserId:        userIdGen.Generate(githubPullRequestComment.AuthorUserId),
				CreatedDate:   githubPullRequestComment.GithubCreatedAt,
				CommitSha:     "",
				Position:      0,
			}
			return []interface{}{
				domainPrComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
