package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertIssueLabels() error {
	var githubIssueLabels []githubModels.GithubIssueLabel

	err := lakeModels.Db.Find(&githubIssueLabels).Error
	if err != nil {
		return err
	}
	domainIdGeneratorIssue := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	for _, githubIssueLabel := range githubIssueLabels {
		domainIl := convertToIssueLabelModel(&githubIssueLabel, domainIdGeneratorIssue)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainIl).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToIssueLabelModel(il *githubModels.GithubIssueLabel,
	domainIdGeneratorIssue *didgen.DomainIdGenerator) *ticket.IssueLabel {
	domainIl := &ticket.IssueLabel{
		IssueId:   domainIdGeneratorIssue.Generate(il.IssueId),
		LabelName: il.LabelName,
	}
	return domainIl
}
