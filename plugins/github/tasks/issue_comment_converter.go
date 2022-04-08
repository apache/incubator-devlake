package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

var ConvertIssueCommentsMeta = core.SubTaskMeta{
	Name:             "convertIssueComments",
	EntryPoint:       ConvertIssueComments,
	EnabledByDefault: true,
	Description:      "ConvertIssueComments data from Github api",
}

func ConvertIssueComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&githubModels.GithubIssueComment{}).
		Joins("left join github_issues "+
			"on github_issues.github_id = github_issue_comments.issue_id").
		Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubIssueComment{}),
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
			githubIssueComment := inputRow.(*githubModels.GithubIssueComment)
			domainIssueComment := &ticket.IssueComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(githubIssueComment.GithubId),
				},
				IssueId:     issueIdGen.Generate(githubIssueComment.IssueId),
				Body:        githubIssueComment.Body,
				UserId:      userIdGen.Generate(githubIssueComment.AuthorUserId),
				CreatedDate: githubIssueComment.GithubCreatedAt,
			}
			return []interface{}{
				domainIssueComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
