package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func PrCommitConvertor(ctx context.Context, repoId int) (err error) {
	githubPullRequestCommit := &models.GithubPullRequestCommit{}

	cursor, err := lakeModels.Db.Model(&githubPullRequestCommit).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_commits.pull_request_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
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
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequestCommit)
		if err != nil {
			return err
		}
		if pullRequestId != githubPullRequestCommit.PullRequestId {
			domainPullRequestId = domainIdGenerator.Generate(githubPullRequestCommit.PullRequestId)
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
