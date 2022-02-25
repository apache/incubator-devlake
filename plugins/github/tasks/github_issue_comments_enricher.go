package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

func EnrichIssueComments(ctx context.Context, repoId int) error {
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
		githubIssueComment.IssueId = issue.GithubId

		err = lakeModels.Db.Save(githubIssueComment).Error
		if err != nil {
			return err
		}
	}
	return nil
}
