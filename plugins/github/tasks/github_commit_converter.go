package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Table("github_commits").
		Where("github_commits._raw_data_params = ?", data.Options.ParamString).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})
	domainRepoId := repoDidGen.Generate(repoId)
	userDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubCommit{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&code.Commit{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
			reflect.TypeOf(&code.RepoCommit{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
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
			domainCommit.RawDataOrigin = githubCommit.RawDataOrigin

			repoCommit := &code.RepoCommit{
				RepoId:    domainRepoId,
				CommitSha: domainCommit.Sha,
			}
			repoCommit.RawDataOrigin = githubCommit.RawDataOrigin

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
