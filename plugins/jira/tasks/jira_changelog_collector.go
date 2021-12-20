package tasks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiAuthor struct {
	Self        string `json:"self,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Active      bool   `json:"active,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
	AccountType string `json:"accountType,omitempty"`
}

type JiraApiChangelogItem struct {
	Field      string `json:"field,omitempty"`
	FieldType  string `json:"fieldType,omitempty"`
	FieldId    string `json:"fieldId,omitempty"`
	From       string `json:"from,omitempty"`
	FromString string `json:"fromString,omitempty"`
	To         string `json:"to,omitempty"`
	ToString   string `json:"toString,omitempty"`
}

type JiraApiChangeLog struct {
	Id      string                 `json:"id,omitempty"`
	Author  JiraApiAuthor          `json:"author,omitempty"`
	Created core.Iso8601Time       `json:"created,omitempty"`
	Items   []JiraApiChangelogItem `json:"items,omitempty"`
}

type JiraApiChangelogsResponse struct {
	JiraPagination
	Values []JiraApiChangeLog `json:"values,omitempty"`
}

func CollectChangelogs(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	ctx context.Context,
) error {
	jiraIssue := &models.JiraIssue{}

	/*
		`CollectIssues` will take into account of `since` option and set the `updated` field for issues that have
		updates, So when it comes to collecting changelogs, we only need to compare an issue's `updated` field with its
		`changelog_updated` field. If `changelog_updated` is older, then we'll collect changelogs for this issue and
		set its `changelog_updated` to `updated` at the end.
	*/
	cursor, err := lakeModels.Db.Model(jiraIssue).
		Select("jira_issues.issue_id", "jira_issues.updated").
		Joins(`LEFT JOIN jira_board_issues ON (
			jira_board_issues.source_id = jira_issues.source_id AND
			jira_board_issues.issue_id = jira_issues.issue_id
		)`).
		Where(`
			jira_board_issues.source_id = ? AND
			jira_board_issues.board_id = ? AND
			(jira_issues.changelog_updated IS NULL OR jira_issues.changelog_updated < jira_issues.updated)
			`,
			source.ID,
			boardId,
		).
		Rows()

	if err != nil {
		return err
	}
	defer cursor.Close()

	changelogScheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer changelogScheduler.Release()
	issueScheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer issueScheduler.Release()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		issueId := jiraIssue.IssueId
		updated := jiraIssue.Updated
		err = issueScheduler.Submit(func() error {
			err = collectChangelogsByIssueId(changelogScheduler, source, jiraApiClient, issueId)
			if err != nil {
				return err
			}
			issue := &models.JiraIssue{SourceId: source.ID, IssueId: issueId}
			err = lakeModels.Db.Model(issue).Update("changelog_updated", updated).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	issueScheduler.WaitUntilFinish()
	changelogScheduler.WaitUntilFinish()

	return nil
}

func collectChangelogsByIssueId(
	scheduler *utils.WorkerScheduler,
	source *models.JiraSource,
	jiraApiClient *JiraApiClient,
	issueId uint64,
) error {
	return jiraApiClient.FetchPages(scheduler, fmt.Sprintf("api/3/issue/%v/changelog", issueId), nil,
		func(res *http.Response) error {
			// parse response
			jiraApiChangelogResponse := &JiraApiChangelogsResponse{}
			err := core.UnmarshalResponse(res, jiraApiChangelogResponse)
			if err != nil {
				return err
			}

			// process changelogs
			for _, jiraApiChangelog := range jiraApiChangelogResponse.Values {

				jiraChangelog, err := convertChangelog(&jiraApiChangelog, source)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				jiraChangelog.IssueId = issueId
				// save changelog
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(jiraChangelog).Error
				if err != nil {
					return err
				}

				// process changelog items
				for _, jiraApiChangelogItem := range jiraApiChangelog.Items {
					jiraChangelogItem, err := convertChangelogItem(
						source,
						jiraChangelog.ChangelogId,
						&jiraApiChangelogItem,
					)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(jiraChangelogItem).Error
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
}

func convertChangelog(jiraApiChangelog *JiraApiChangeLog, source *models.JiraSource) (*models.JiraChangelog, error) {
	id, err := strconv.ParseUint(jiraApiChangelog.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	return &models.JiraChangelog{
		SourceId:          source.ID,
		ChangelogId:       id,
		AuthorAccountId:   jiraApiChangelog.Author.AccountId,
		AuthorDisplayName: jiraApiChangelog.Author.DisplayName,
		AuthorActive:      jiraApiChangelog.Author.Active,
		Created:           jiraApiChangelog.Created.ToTime(),
	}, nil
}

func convertChangelogItem(
	source *models.JiraSource,
	changelogId uint64,
	jiraApiChangeItem *JiraApiChangelogItem,
) (*models.JiraChangelogItem, error) {
	return &models.JiraChangelogItem{
		SourceId:    source.ID,
		ChangelogId: changelogId,
		Field:       jiraApiChangeItem.Field,
		FieldType:   jiraApiChangeItem.FieldType,
		FieldId:     jiraApiChangeItem.FieldId,
		From:        jiraApiChangeItem.From,
		FromString:  jiraApiChangeItem.FromString,
		To:          jiraApiChangeItem.To,
		ToString:    jiraApiChangeItem.ToString,
	}, nil
}
