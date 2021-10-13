package tasks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func GetWhereClauseConditionally(latestUpdatedIssue models.JiraIssue, since time.Time) string {
	var whereClause string

	if latestUpdatedIssue.ID > 0 {
		// This is not the first time we have fetched data for Jira.
		// Therefore only get data since the last time we fetched data
		whereClause = fmt.Sprintf(`jira_board_issues.board_id = ?
		AND (jira_issues.changelog_updated is null OR '%v' < jira_issues.updated)`, latestUpdatedIssue.Updated)
	} else if !since.IsZero() {
		// This is the first time we have fetched data
		// "Since" was provided in the POST request so we start there
		whereClause = fmt.Sprintf(`jira_board_issues.board_id = ?
		AND (jira_issues.changelog_updated is null OR '%v' < jira_issues.updated)`, since)
	} else {
		// This is the first time we fetch the data and since was not provided
		whereClause = "jira_board_issues.board_id = ?"
	}
	return whereClause
}

func GetLatestIssueFromDB() models.JiraIssue {
	var latestUpdatedIssue models.JiraIssue
	err := lakeModels.Db.Debug().Order("changelog_updated DESC").Limit(1).Find(&latestUpdatedIssue).Error
	if err != nil {
		logger.Error("err", err)
	}
	return latestUpdatedIssue
}

func CollectChangelogs(boardId uint64, since time.Time, ctx context.Context) error {
	jiraIssue := &models.JiraIssue{}

	// Get "Latest Issue" from the DB
	latestUpdatedIssue := GetLatestIssueFromDB()

	whereClause := GetWhereClauseConditionally(latestUpdatedIssue, since)

	// Get all Issues from 'changelog_updated' time on latest Issue.
	// Then get Changelogs for those issues.

	cursor, err := lakeModels.Db.Debug().Model(jiraIssue).
		Select("jira_issues.id", "jira_issues.updated").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.id").
		Where(whereClause,
			boardId).
		Rows()

	if err != nil {
		return err
	}
	defer cursor.Close()

	changelogScheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	issueScheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer changelogScheduler.Release()
	defer issueScheduler.Release()
	jiraApiClient := GetJiraApiClient()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		issueId := jiraIssue.ID
		updated := jiraIssue.Updated
		err = issueScheduler.Submit(func() error {
			err = collectChangelogsByIssueId(changelogScheduler, jiraApiClient, issueId)
			if err != nil {
				return err
			}
			issue := &models.JiraIssue{Model: lakeModels.Model{ID: issueId}}
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

func collectChangelogsByIssueId(scheduler *utils.WorkerScheduler, jiraApiClient *JiraApiClient, issueId uint64) error {
	return jiraApiClient.FetchPages(scheduler, fmt.Sprintf("/api/3/issue/%v/changelog", issueId), nil,
		func(res *http.Response) error {
			// parse response
			jiraApiChangelogResponse := &JiraApiChangelogsResponse{}
			err := core.UnmarshalResponse(res, jiraApiChangelogResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			// process changelogs
			for _, jiraApiChangelog := range jiraApiChangelogResponse.Values {

				jiraChangelog, err := convertChangelog(&jiraApiChangelog)
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
					logger.Error("Error: ", err)
					return err
				}

				// process changelog items
				lakeModels.Db.Delete(models.JiraChangelogItem{}, "changelog_id = ?", jiraChangelog.ID)
				for _, jiraApiChangelogItem := range jiraApiChangelog.Items {
					jiraChangelogItem, err := convertChangelogItem(jiraChangelog.ID, &jiraApiChangelogItem)
					if err != nil {
						logger.Error("Error: ", err)
						return err
					}
					// save changelog item
					err = lakeModels.Db.Create(jiraChangelogItem).Error
					if err != nil {
						logger.Error("Error: ", err)
						return err
					}
				}
			}
			return nil
		})
}

func convertChangelog(jiraApiChangelog *JiraApiChangeLog) (*models.JiraChangelog, error) {
	id, err := strconv.ParseUint(jiraApiChangelog.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	return &models.JiraChangelog{
		Model:             lakeModels.Model{ID: id},
		AuthorAccountId:   jiraApiChangelog.Author.AccountId,
		AuthorDisplayName: jiraApiChangelog.Author.DisplayName,
		AuthorActive:      jiraApiChangelog.Author.Active,
		Created:           jiraApiChangelog.Created.ToTime(),
	}, nil
}

func convertChangelogItem(changelogId uint64, jiraApiChangeItem *JiraApiChangelogItem) (*models.JiraChangelogItem, error) {
	return &models.JiraChangelogItem{
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
