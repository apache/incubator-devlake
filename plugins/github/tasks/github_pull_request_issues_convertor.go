package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequestIssues(ctx context.Context, repoId int) error {
	githubPullRequestIssue := &githubModels.GithubPullRequestIssue{}
	cursor, err := lakeModels.Db.Model(githubPullRequestIssue).
		Select("github_pull_request_issues.*").
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
			err := lakeModels.Db.Where("pull_request_id = ?",
				pullRequestId).Delete(&code.PullRequestLabel{}).Error
			if err != nil {
				return err
			}
			lastPrId = githubPullRequestIssue.PullRequestId
		}

		pullRequestIssue := &code.PullRequestIssue{
			PullRequestId: pullRequestId,
			IssueId:       issueId,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(pullRequestIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
