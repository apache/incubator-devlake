package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequestIssues(ctx context.Context, repoId int) error {
	githubPullRequestIssue := &githubModels.GithubPullRequestIssue{}
	cursor, err := lakeModels.Db.Model(githubPullRequestIssue).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_issues.pull_request_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
		Order("pull_request_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorPr := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})
	domainIdGeneratorIssue := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	lastPrId := 0
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubPullRequestIssue)
		if err != nil {
			return err
		}
		pullRequestId := domainIdGeneratorPr.Generate(githubPullRequestIssue.PullRequestId)
		issueId := domainIdGeneratorIssue.Generate(githubPullRequestIssue.IssueId)
		if lastPrId != githubPullRequestIssue.PullRequestId {
			// Clean up old data
			err = lakeModels.Db.Where("pull_request_id = ?",
				pullRequestId).Delete(&crossdomain.PullRequestIssue{}).Error
			if err != nil {
				return err
			}
			lastPrId = githubPullRequestIssue.PullRequestId
		}

		pullRequestIssue := &crossdomain.PullRequestIssue{
			PullRequestId: pullRequestId,
			IssueId:       issueId,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(pullRequestIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
