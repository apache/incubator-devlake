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

func ConvertMrs(ctx context.Context, projectId int) error {
	domainMrIdGenerator := didgen.NewDomainIdGenerator(&gitlabModels.GitlabMergeRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&gitlabModels.GitlabProject{})

	gitlabMr := &gitlabModels.GitlabMergeRequest{}
	//Find all piplines associated with the current projectid
	cursor, err := lakeModels.Db.Model(gitlabMr).Where("project_id=?", projectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, gitlabMr)

		domainPr := &code.PullRequest{
			DomainEntity: domainlayer.DomainEntity{
				Id: domainMrIdGenerator.Generate(gitlabMr.GitlabId),
			},
			RepoId:      domainRepoIdGenerator.Generate(gitlabMr.ProjectId),
			Status:      gitlabMr.State,
			Title:       gitlabMr.Title,
			Url:         gitlabMr.WebUrl,
			CreatedDate: gitlabMr.GitlabCreatedAt,
			MergedDate:  gitlabMr.MergedAt,
			ClosedAt:    gitlabMr.ClosedAt,
		}
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainPr).Error
		if err != nil {
			return err
		}
	}
	return nil
}
