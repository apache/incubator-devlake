package tasks

import (
	"context"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertIssueLabels(ctx context.Context) error {
	githubIssueLabel := &githubModels.GithubIssueLabel{}
	cursor, err := lakeModels.Db.Model(githubIssueLabel).
		Select("github_issue_labels.*").
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
			return core.TaskCanceled
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
