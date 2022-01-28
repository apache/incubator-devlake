package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertReposToBoard() error {
	var githubRepositorys []githubModels.GithubRepository
	err := lakeModels.Db.Find(&githubRepositorys).Error
	if err != nil {
		return err
	}
	domainIdGenerator := didgen.NewDomainIdGenerator(&githubModels.GithubRepository{})
	for _, repository := range githubRepositorys {
		domainBoard := convertToBoardModel(&repository, domainIdGenerator)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainBoard).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToBoardModel(repository *githubModels.GithubRepository, domainIdGenerator *didgen.DomainIdGenerator) *ticket.Board {
	domainBoard := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: domainIdGenerator.Generate(repository.GithubId),
		},
		Name:        repository.Name,
		Url:         repository.HTMLUrl,
		Description: repository.Description,

		CreatedDate: &repository.CreatedDate,
	}
	return domainBoard
}
