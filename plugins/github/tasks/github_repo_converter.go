package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertRepos(ctx context.Context) error {
	githubRepository := &githubModels.GithubRepo{}
	cursor, err := lakeModels.Db.Model(githubRepository).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	//Will be used when generating domainId, to avoid to compile every iteration
	domainUserIdGenerator := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubRepository)
		if err != nil {
			return err
		}
		domainRepository := convertToRepositoryModel(githubRepository, domainRepoIdGenerator)
		domainRepository.OwnerId = domainUserIdGenerator.Generate(githubRepository.OwnerId)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainRepository).Error
		if err != nil {
			return err
		}
		domainBoard := convertToBoardModel(githubRepository, domainRepoIdGenerator)
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainBoard).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func convertToRepositoryModel(repository *githubModels.GithubRepo, domainIdGenerator *didgen.DomainIdGenerator) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: domainIdGenerator.Generate(repository.GithubId),
		},
		Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
		Url:         repository.HTMLUrl,
		Description: repository.Description,
		ForkedFrom:  repository.ParentHTMLUrl,
		Language:    repository.Language,
		CreatedDate: repository.CreatedDate,
		UpdatedDate: repository.UpdatedDate,
	}
	return domainRepository
}

func convertToBoardModel(repository *githubModels.GithubRepo, domainIdGenerator *didgen.DomainIdGenerator) *ticket.Board {
	domainBoard := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: domainIdGenerator.Generate(repository.GithubId),
		},
		Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
		Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
		Description: repository.Description,

		CreatedDate: &repository.CreatedDate,
	}
	return domainBoard
}
