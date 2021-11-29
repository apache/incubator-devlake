package tasks

import (
	"database/sql"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/base"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertChangelogs(sourceId uint64, boardId uint64) error {

	jiraChangelog := &jiraModels.JiraChangelog{}

	var c1, c2 *sql.Rows
	var err error
	defer func() {
		if c1 != nil {
			c1.Close()
		}
		if c2 != nil {
			c2.Close()
		}
	}()
	// select all changelogs belongs to the board
	c1, err = lakeModels.Db.Model(jiraChangelog).
		Select("jira_changelogs.*").
		Joins(`left join jira_board_issues on (jira_board_issues.issue_id = jira_changelogs.issue_id)`).
		Where("jira_board_issues.source_id = ? AND jira_board_issues.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}

	issueOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraIssue{})
	changelogOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraChangelogItem{})

	// iterate all rows
	for c1.Next() {
		err = lakeModels.Db.ScanRows(c1, jiraChangelog)
		if err != nil {
			return err
		}

		var items []jiraModels.JiraChangelogItem
		err = lakeModels.Db.Where("changelog_id = ?", jiraChangelog.ChangelogId).Find(&items).Error
		if err != nil {
			return err
		}
		for _, jiraChangelogItem := range items {
			if err != nil {
				return err
			}
			changelog := &ticket.Changelog{
				DomainEntity: base.DomainEntity{
					OriginKey: changelogOriginKeyGenerator.Generate(
						jiraChangelogItem.SourceId,
						jiraChangelogItem.ChangelogId,
						jiraChangelogItem.Field,
					),
				},
				IssueOriginKey: issueOriginKeyGenerator.Generate(jiraChangelog.SourceId, jiraChangelog.IssueId),
				AuthorName:     jiraChangelog.AuthorDisplayName,
				FieldName:      jiraChangelogItem.Field,
				From:           jiraChangelogItem.FromString,
				To:             jiraChangelogItem.ToString,
				CreatedDate:    jiraChangelog.Created,
			}

			err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(changelog).Error
			if err != nil {
				println("err", err)
				return err
			}
		}
	}
	return nil
}
