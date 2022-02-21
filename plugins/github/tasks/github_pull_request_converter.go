package tasks

import (
	"context"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequests(ctx context.Context) error {
	githubPullRequest := &githubModels.GithubPullRequest{}
	cursor, err := lakeModels.Db.Model(githubPullRequest).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainPrIdGenerator := didgen.NewDomainIdGenerator(githubPullRequest)

	for cursor.Next() {
		select {
		case <-ctx.Done():
			return core.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequest)
		if err != nil {
			return err
		}
		domainPr := convertToPullRequestModel(githubPullRequest, domainPrIdGenerator)
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainPr).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToPullRequestModel(pr *githubModels.GithubPullRequest, domainGenerator *didgen.DomainIdGenerator) *code.PullRequest {
	domainPr := &code.PullRequest{
		DomainEntity: domainlayer.DomainEntity{
			Id: domainGenerator.Generate(pr.GithubId),
		},
		RepoId:         uint64(pr.RepoId),
		Status:         pr.State,
		Title:          pr.Title,
		CreatedDate:    pr.GithubCreatedAt,
		MergedDate:     pr.MergedAt,
		ClosedAt:       pr.ClosedAt,
		Type:           pr.Type,
		Component:      pr.Component,
		MergeCommitSha: pr.MergeCommitSha,
		BaseRef:        pr.BaseRef,
		BaseCommitSha:  pr.BaseCommitSha,
		HeadRef:        pr.HeadRef,
		HeadCommitSha:  pr.HeadCommitSha,
	}
	return domainPr
}
