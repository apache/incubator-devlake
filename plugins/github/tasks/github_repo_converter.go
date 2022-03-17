package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertRepos(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&githubModels.GithubRepo{}).
		Where("github_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubRepo{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_REPOSITORIES_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			repository := inputRow.(*githubModels.GithubRepo)
			domainRepository := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(repository.GithubId),
				},
				Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
				Url:         repository.HTMLUrl,
				Description: repository.Description,
				ForkedFrom:  repository.ParentHTMLUrl,
				Language:    repository.Language,
				CreatedDate: repository.CreatedDate,
				UpdatedDate: repository.UpdatedDate,
			}

			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(repository.GithubId),
				},
				Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
				Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
				Description: repository.Description,
				CreatedDate: &repository.CreatedDate,
			}

			return []interface{}{
				domainRepository,
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
