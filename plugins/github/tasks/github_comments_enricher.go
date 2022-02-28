package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

//EnrichComments will enrich issue comments with issue number
//also move pr comments to github_pull_request_comments
func EnrichComments(ctx context.Context, repoId int) error {
	cursor, err := lakeModels.Db.Model(&githubModels.GithubIssueComment{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		githubIssueComment := &githubModels.GithubIssueComment{}

		err = lakeModels.Db.ScanRows(cursor, githubIssueComment)
		if err != nil {
			return err
		}

		issue := &githubModels.GithubIssue{}
		err = lakeModels.Db.Where("number = ? and `repo_id` = ?", githubIssueComment.IssueNumber, repoId).Limit(1).Find(issue).Error
		if err != nil {
			return err
		}
		//if we can not find issues with issue number above, move the comments to github_pull_request_comments
		if issue.GithubId == 0 {
			pr := &githubModels.GithubPullRequest{}
			err = lakeModels.Db.Where("number = ? and `repo_id` = ?", githubIssueComment.IssueNumber, repoId).Limit(1).Find(pr).Error
			if err != nil {
				return err
			}
			githubPrComment := &githubModels.GithubPullRequestComment{
				GithubId:        githubIssueComment.GithubId,
				PullRequestId:   pr.GithubId,
				Body:            githubIssueComment.Body,
				AuthorUsername:  githubIssueComment.AuthorUsername,
				GithubCreatedAt: githubIssueComment.GithubCreatedAt,
			}
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(githubPrComment).Error
			if err != nil {
				return err
			}
			err = lakeModels.Db.Delete(githubIssueComment).Error
			if err != nil {
				return err
			}
			continue
		}
		githubIssueComment.IssueId = issue.GithubId

		err = lakeModels.Db.Save(githubIssueComment).Error
		if err != nil {
			return err
		}
	}
	return nil
}
