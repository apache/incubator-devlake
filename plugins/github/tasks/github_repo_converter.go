package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertRepos() error {
	var githubRepositorys []githubModels.GithubRepo
	err := lakeModels.Db.Find(&githubRepositorys).Error
	if err != nil {
		return err
	}
	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	for _, repository := range githubRepositorys {
		domainRepository := convertToRepositoryModel(&repository)
		domainRepository.OwnerId = userIdGen.Generate(repository.OwnerId)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainRepository).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToRepositoryModel(repository *githubModels.GithubRepo) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(repository).Generate(repository.GithubId),
		},
		Name:        repository.Name,
		Url:         repository.HTMLUrl,
		Description: repository.Description,
		ForkedFrom:  repository.ParentHTMLUrl,
		Language:    repository.Language,
		CreatedDate: repository.CreatedDate,
		UpdatedDate: repository.UpdatedDate,
	}
	return domainRepository
}
