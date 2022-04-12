package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_commits into  domain layer table commits",
}

func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Table("_tool_github_commits gc").
		Joins(`left join _tool_github_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`).
		Select("gc.*").
		Where("grc.repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})
	domainRepoId := repoDidGen.Generate(repoId)
	userDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		InputRowType: reflect.TypeOf(githubModels.GithubCommit{}),
		Input:        cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubCommit := inputRow.(*githubModels.GithubCommit)
			domainCommit := &code.Commit{
				Sha:            githubCommit.Sha,
				Message:        githubCommit.Message,
				Additions:      githubCommit.Additions,
				Deletions:      githubCommit.Deletions,
				AuthorId:       userDidGen.Generate(githubCommit.AuthorId),
				AuthorName:     githubCommit.AuthorName,
				AuthorEmail:    githubCommit.AuthorEmail,
				AuthoredDate:   githubCommit.AuthoredDate,
				CommitterName:  githubCommit.CommitterName,
				CommitterEmail: githubCommit.CommitterEmail,
				CommittedDate:  githubCommit.CommittedDate,
				CommitterId:    userDidGen.Generate(githubCommit.CommitterId),
			}
			repoCommit := &code.RepoCommit{
				RepoId:    domainRepoId,
				CommitSha: domainCommit.Sha,
			}
			return []interface{}{
				domainCommit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
