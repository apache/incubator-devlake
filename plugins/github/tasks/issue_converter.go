package tasks

import (
	"fmt"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
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

func convertStateToStatus(state string) string {
	if state == "closed" {
		return "Resolved"
	} else {
		return "Todo"
	}
}

func convertToIssueModel(issue *githubModels.GithubIssue) *ticket.Issue {
	domainIssue := &ticket.Issue{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(issue).Generate(issue.GithubId),
		},
		Key:               fmt.Sprint(issue.GithubId),
		Title:             issue.Title,
		Summary:           issue.Body,
		Status:            convertStateToStatus(issue.State),
		Priority:          issue.Priority,
		Type:              issue.Type,
		AssigneeOriginKey: issue.Assignee,
		LeadTimeMinutes:   issue.LeadTimeMinutes,
		CreatedDate:       issue.GithubCreatedAt,
		UpdatedDate:       issue.GithubUpdatedAt,
		ResolutionDate:    issue.ClosedAt,
	}
	return domainIssue
}
