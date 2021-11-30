package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertRepos() error {
	var githubRepositorys []githubModels.GithubRepository
	err := lakeModels.Db.Find(&githubRepositorys).Error
	if err != nil {
		return err
	}
	for _, repository := range githubRepositorys {
		domainRepository := convertToRepositoryModel(&repository)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainRepository).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToRepositoryModel(repository *githubModels.GithubRepository) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(repository).Generate(repository.GithubId),
		},
		Name: repository.Name,
		Url:  repository.HTMLUrl,
	}
	return domainRepository
}
