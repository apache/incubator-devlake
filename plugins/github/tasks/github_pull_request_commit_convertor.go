package tasks

import (
	"context"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func PrCommitConvertor(ctx context.Context) (err error) {
	githubPullRequestCommit := &models.GithubPullRequestCommit{}

	cursor, err := lakeModels.Db.Model(&githubPullRequestCommit).
		Order("pull_request_id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	var pullRequestId int
	domainPullRequestId := ""
	domainIdGenerator := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return core.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequestCommit)
		if err != nil {
			return err
		}
		if pullRequestId != githubPullRequestCommit.PullRequestId {
			domainPullRequestId = domainIdGenerator.Generate(pullRequestId)
			err := lakeModels.Db.Where("pull_request_id = ?",
				domainPullRequestId).Delete(&code.PullRequestCommit{}).Error
			if err != nil {
				return err
			}
			pullRequestId = githubPullRequestCommit.PullRequestId
		}

		if err != nil {
			return err
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&code.PullRequestCommit{
			CommitSha:     githubPullRequestCommit.CommitSha,
			PullRequestId: domainPullRequestId,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
