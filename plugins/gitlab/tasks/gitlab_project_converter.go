package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertProjects(ctx context.Context, projectId int) error {

	gitlabProject := &gitlabModels.GitlabProject{}
	//Find all piplines associated with the current projectid
	cursor, err := lakeModels.Db.Model(gitlabProject).Where("gitlab_id=?", projectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		domainRepository := convertToRepositoryModel(gitlabProject)
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
			Id: didgen.NewDomainIdGenerator(project).Generate(project.GitlabId),
		},
		Name:        project.Name,
		Url:         project.WebUrl,
		Description: project.Description,
		ForkedFrom:  project.ForkedFromProjectWebUrl,
		CreatedDate: project.CreatedDate,
		UpdatedDate: project.UpdatedDate,
	}
	return domainRepository
}
