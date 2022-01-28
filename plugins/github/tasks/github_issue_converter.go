package tasks

import (
	"fmt"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertIssues(repoId int) error {
	var githubIssues []githubModels.GithubIssue
	err := lakeModels.Db.Find(&githubIssues).Error
	if err != nil {
		return err
	}
	domainIdGeneratorIssue := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	domainIdGeneratorGithubUser := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	boardIssue := &ticket.BoardIssue{
		BoardId: didgen.NewDomainIdGenerator(&githubModels.GithubRepository{}).Generate(repoId),
	}
	for _, issue := range githubIssues {
		domainIssue := convertToIssueModel(&issue, domainIdGeneratorIssue, domainIdGeneratorGithubUser)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainIssue).Error
		if err != nil {
			return err
		}
		boardIssue.IssueId = domainIssue.Id

		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(boardIssue).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func convertStateToStatus(state string) string {
	if state == "closed" {
		return ticket.DONE
	} else {
		return ticket.TODO
	}
}

func convertToIssueModel(issue *githubModels.GithubIssue, domainIdGeneratorIssue *didgen.DomainIdGenerator,
	domainIdGeneratorGithubUser *didgen.DomainIdGenerator) *ticket.Issue {
	domainIssue := &ticket.Issue{
		DomainEntity:    domainlayer.DomainEntity{Id: domainIdGeneratorIssue.Generate(issue.GithubId)},
		Key:             fmt.Sprint(issue.GithubId),
		Title:           issue.Title,
		Summary:         issue.Body,
		Status:          convertStateToStatus(issue.State),
		Priority:        issue.Priority,
		Type:            issue.Type,
		AssigneeId:      domainIdGeneratorGithubUser.Generate(issue.AssigneeId),
		AssigneeName:    issue.AssigneeName,
		LeadTimeMinutes: issue.LeadTimeMinutes,
		CreatedDate:     &issue.GithubCreatedAt,
		UpdatedDate:     &issue.GithubUpdatedAt,
		ResolutionDate:  issue.ClosedAt,
		Severity:        issue.Severity,
		Component:       issue.Component,
	}
	return domainIssue
}
