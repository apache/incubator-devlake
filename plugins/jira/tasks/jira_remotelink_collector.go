package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

var ErrNotFoundResource = errors.New("not found the resource")

type JiraApiRemotelink struct {
	Id           uint64
	Self         string
	Application  json.RawMessage
	Relationship string
	Object       struct {
		Url    string
		Title  string
		Icon   json.RawMessage
		Status json.RawMessage
	}
}

// need to store a origin json body into RawJson, by this approach, we dont need to Marshal it back to bytes
type JiraApiRemotelinksResponse []json.RawMessage
type remoteLinkCollector func(
	source *models.JiraSource,
	jiraApiClient *JiraApiClient,
	issueId uint64,
) error

func CollectRemoteLinks(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	rateLimitPerSecondInt int,
	ctx context.Context,
	collector remoteLinkCollector,
) error {
	jiraIssue := &models.JiraIssue{}

	/*
		`CollectIssues` will take into account of `since` option and set the `updated` field for issues that have
		updates, So when it comes to collecting remotelinks, we only need to compare an issue's `updated` field with its
		`remotelink_updated` field. If `remotelink_updated` is older, then we'll collect remotelinks for this issue and
		set its `remotelink_updated` to `updated` at the end.
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
			(jira_issues.remotelink_updated IS NULL OR jira_issues.remotelink_updated < jira_issues.updated)
			`,
			source.ID,
			boardId,
		).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueScheduler, err := utils.NewWorkerScheduler(10, rateLimitPerSecondInt, ctx)
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
			err = collector(source, jiraApiClient, issueId)
			if err == ErrNotFoundResource {
				return nil
			}
			if err != nil {
				return err
			}
			issue := &models.JiraIssue{SourceId: source.ID, IssueId: issueId}
			err = lakeModels.Db.Model(issue).Update("remotelink_updated", updated).Error
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

	return nil
}

func collectRemotelinksByIssueId(
	source *models.JiraSource,
	jiraApiClient *JiraApiClient,
	issueId uint64,
) error {
	res, err := jiraApiClient.Get(fmt.Sprintf("api/3/issue/%v/remotelink", issueId), nil, nil)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return ErrNotFoundResource
	}
	apiRemotelinks := &JiraApiRemotelinksResponse{}
	err = core.UnmarshalResponse(res, apiRemotelinks)
	if err != nil {
		return err
	}

	apiRemotelink := &JiraApiRemotelink{}
	remotelink := &models.JiraRemotelink{}

	// delete previous collected remotelink
	err = lakeModels.Db.Where("source_id = ? AND issue_id = ?", source.ID, issueId).Delete(remotelink).Error
	if err != nil {
		return err
	}

	for _, apiRemotelinkRaw := range *apiRemotelinks {
		// unmarshal to fetch id for primary key
		err = json.Unmarshal(apiRemotelinkRaw, apiRemotelink)
		if err != nil {
			return err
		}
		// create a empty record with pk only
		remotelink.SourceId = source.ID
		remotelink.IssueId = issueId
		remotelink.RemotelinkId = apiRemotelink.Id
		remotelink.RawJson = datatypes.JSON(apiRemotelinkRaw)
		// save raw response, delay feilds extraction to enrich stage
		err = lakeModels.Db.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"raw_json"}),
		}).Create(remotelink).Error
		if err != nil {
			return err
		}
	}
	return nil
}
