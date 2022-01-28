package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequests() error {
	var githubPullRequests []githubModels.GithubPullRequest
	err := lakeModels.Db.Find(&githubPullRequests).Error
	if err != nil {
		return err
	}
	for _, pullrequest := range githubPullRequests {
		domainPr := convertToPullRequestModel(&pullrequest)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainPr).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToPullRequestModel(pr *githubModels.GithubPullRequest) *code.PullRequest {
	domainPr := &code.PullRequest{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(pr).Generate(pr.GithubId),
		},
		RepoId:         uint64(pr.RepositoryId),
		Status:         pr.State,
		Title:          pr.Title,
		CreatedDate:    pr.GithubCreatedAt,
		MergedDate:     pr.MergedAt,
		ClosedAt:       pr.ClosedAt,
		Type:           pr.Type,
		Component:      pr.Component,
		MergeCommitSha: pr.MergeCommitSha,
	}
	return domainPr
}
