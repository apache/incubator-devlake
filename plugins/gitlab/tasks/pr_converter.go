package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertPrs() error {
	var gitlabMrs []gitlabModels.GitlabMergeRequest
	err := lakeModels.Db.Find(&gitlabMrs).Error
	if err != nil {
		return err
	}
	for _, mr := range gitlabMrs {
		domainPr := convertToPrModel(&mr)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainPr).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToPrModel(mr *gitlabModels.GitlabMergeRequest) *code.Pr {
	domainPr := &code.Pr{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(mr).Generate(mr.GitlabId),
		},
		RepoId:      uint64(mr.ProjectId),
		State:       mr.State,
		Title:       mr.Title,
		Url:         mr.WebUrl,
		CreatedDate: mr.GitlabCreatedAt,
		MergedDate:  mr.MergedAt,
		ClosedAt:    mr.ClosedAt,
	}
	return domainPr
}
