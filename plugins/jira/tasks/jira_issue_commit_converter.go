package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertIssueCommits(sourceId uint64, boardId uint64) error {
	// select all changelogs belongs to the board
	cursor, err := lakeModels.Db.Table("jira_issue_commits jic").
		Joins(`left join jira_board_issues jbi on (
			jbi.source_id = jic.source_id
			AND jbi.issue_id = jic.issue_id
		)`).
		Select("jic.*").
		Where("jbi.source_id = ? AND jbi.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&jiraModels.JiraIssue{})

	// save in batch
	size := 1000
	i := 0
	batch := make([]crossdomain.IssueCommit, size)
	saveBatch := func() error {
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(batch[:i], size).Error
		if err != nil {
			return err
		}
		return nil
	}
	row := &jiraModels.JiraIssueCommit{}
	// iterate all rows
	for cursor.Next() {
		if i >= size {
			err = saveBatch()
			if err != nil {
				return err
			}
			i = 0
		}
		err = lakeModels.Db.ScanRows(cursor, row)
		if err != nil {
			return err
		}
		issueCommit := &batch[i]
		issueCommit.IssueId = issueIdGenerator.Generate(row.SourceId, row.IssueId)
		issueCommit.CommitSha = row.CommitSha
		i++
	}
	if i > 0 {
		err = saveBatch()
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}
