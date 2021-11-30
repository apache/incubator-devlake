package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertRepos() error {
	var gitlabProjects []gitlabModels.GitlabProject
	err := lakeModels.Db.Find(&gitlabProjects).Error
	if err != nil {
		return err
	}
	for _, repository := range gitlabProjects {
		domainRepository := convertToRepositoryModel(&repository)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainRepository).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToRepositoryModel(project *gitlabModels.GitlabProject) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(project).Generate(project.GitlabId),
		},
		Name: project.Name,
		Url:  project.WebUrl,
	}
	return domainRepository
}
