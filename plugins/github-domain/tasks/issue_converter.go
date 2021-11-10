package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
	"github.com/merico-dev/lake/plugins/domainlayer/models/ticket"
	"github.com/merico-dev/lake/plugins/domainlayer/okgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertIssues() error {
	var githubIssues []githubModels.GithubIssue
	err := lakeModels.Db.Find(&githubIssues).Error
	if err != nil {
		return err
	}
	for _, issue := range githubIssues {
		domainIssue := convertToIssueModel(&issue)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainIssue).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func convertToIssueModel(issue *githubModels.GithubIssue) *ticket.Issue {
	domainIssue := &ticket.Issue{
		DomainEntity: base.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(issue).Generate(issue.GithubId),
		},
		Key:            issue.Title,
		Summary:        issue.Body,
		Status:         issue.State,
		Priority:       issue.Priority,
		Type:           issue.Type,
		CreatedDate:    issue.GithubCreatedAt,
		UpdatedDate:    issue.GithubUpdatedAt,
		ResolutionDate: issue.ClosedAt,
	}
	return domainIssue
}
