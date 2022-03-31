package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequests(ctx context.Context, repoId int) error {
	pr := &githubModels.GithubPullRequest{}
	cursor, err := lakeModels.Db.Model(pr).Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainPrIdGenerator := didgen.NewDomainIdGenerator(pr)
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})

	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, pr)
		if err != nil {
			return err
		}
		domainPr := &code.PullRequest{
			DomainEntity: domainlayer.DomainEntity{
				Id: domainPrIdGenerator.Generate(pr.GithubId),
			},
			RepoId:         domainRepoIdGenerator.Generate(pr.RepoId),
			Status:         pr.State,
			Title:          pr.Title,
			Description:    pr.Body,
			Url:            pr.Url,
			AuthorName:     pr.AuthorName,
			AuthorId:       pr.AuthorId,
			CreatedDate:    pr.GithubCreatedAt,
			MergedDate:     pr.MergedAt,
			ClosedAt:       pr.ClosedAt,
			Key:            pr.Number,
			Type:           pr.Type,
			Component:      pr.Component,
			MergeCommitSha: pr.MergeCommitSha,
			BaseRef:        pr.BaseRef,
			BaseCommitSha:  pr.BaseCommitSha,
			HeadRef:        pr.HeadRef,
			HeadCommitSha:  pr.HeadCommitSha,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainPr).Error
		if err != nil {
			return err
		}
	}
	return nil
}
