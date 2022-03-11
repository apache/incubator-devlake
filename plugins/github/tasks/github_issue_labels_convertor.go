package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertIssueLabels(ctx context.Context, repoId int) error {
	githubIssueLabel := &githubModels.GithubIssueLabel{}
	cursor, err := lakeModels.Db.Model(githubIssueLabel).
		Joins(`left join github_issues on github_issues.github_id = github_issue_labels.issue_id`).
		Where("github_issues.repo_id = ?", repoId).
		Order("issue_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorIssue := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	lastIssueId := 0
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubIssueLabel)
		if err != nil {
			return err
		}
		issueId := domainIdGeneratorIssue.Generate(githubIssueLabel.IssueId)
		if lastIssueId != githubIssueLabel.IssueId {
			// Clean up old data
			err := lakeModels.Db.Where("issue_id = ?",
				issueId).Delete(&ticket.IssueLabel{}).Error
			if err != nil {
				return err
			}
			lastIssueId = githubIssueLabel.IssueId
		}

		issueLabel := &ticket.IssueLabel{
			IssueId:   issueId,
			LabelName: githubIssueLabel.LabelName,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(issueLabel).Error
		if err != nil {
			return err
		}
	}
	return nil
}
