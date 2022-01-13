package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type ChangelogItemResult struct {
	SourceId          uint64 `gorm:"primaryKey"`
	ChangelogId       uint64 `gorm:"primaryKey"`
	Field             string `gorm:"primaryKey"`
	FieldType         string
	FieldId           string
	From              string
	FromString        string
	To                string
	ToString          string
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	Created           time.Time
}

func ConvertChangelogs(sourceId uint64, boardId uint64) error {
	// select all changelogs belongs to the board
	cursor, err := lakeModels.Db.Table("jira_changelog_items").
		Joins(`left join jira_changelogs on (
			jira_changelogs.source_id = jira_changelog_items.source_id
			AND jira_changelogs.changelog_id = jira_changelog_items.changelog_id
		)`).
		Joins(`left join jira_board_issues on (
			jira_board_issues.source_id = jira_changelogs.source_id
			AND jira_board_issues.issue_id = jira_changelogs.issue_id
		)`).
		Select("jira_changelog_items.*, jira_changelogs.*").
		Where("jira_changelog_items.source_id = ? AND jira_board_issues.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&jiraModels.JiraIssue{})
	changelogIdGenerator := didgen.NewDomainIdGenerator(&jiraModels.JiraChangelogItem{})

	// save in batch
	size := 1000
	i := 0
	batch := make([]ticket.Changelog, size)
	saveBatch := func() error {
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(batch[:i], size).Error
		if err != nil {
			println("err", err)
			return err
		}
		logger.Info("convert changelog", fmt.Sprintf("%s .. %s", batch[0].Id, batch[i-1].Id))
		return nil
	}
	sprintIssueConverter := NewSprintIssueConverter()
	row := &ChangelogItemResult{}
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
		changelog := &batch[i]
		changelog.DomainEntity.Id = changelogIdGenerator.Generate(
			row.SourceId,
			row.ChangelogId,
			row.Field,
		)
		changelog.IssueId = issueIdGenerator.Generate(row.SourceId, row.IssueId)
		changelog.AuthorName = row.AuthorDisplayName
		changelog.FieldName = row.Field
		changelog.From = row.FromString
		changelog.To = row.ToString
		changelog.CreatedDate = row.Created
		i++
		sprintIssueConverter.FeedIn(sourceId, *row)
	}
	if i > 0 {
		err = saveBatch()
		if err != nil {
			return err
		}
	}
	err = sprintIssueConverter.UpdateSprintIssue()
	if err != nil {
		logger.Error("update sprint issue error:", err)
	}
	return nil
}
