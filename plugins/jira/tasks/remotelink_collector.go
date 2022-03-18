package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"io/ioutil"
	"net/http"
)

const RAW_REMOTELINK_TABLE = "jira_api_remotelinks"

var _ core.SubTaskEntryPoint = CollectApiRemotelinks

func CollectApiRemotelinks(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	logger := taskCtx.GetLogger()
	logger.Info("collect remotelink")
	jiraIssue := &models.JiraIssue{}

	/*
		`CollectIssues` will take into account of `since` option and set the `updated` field for issues that have
		updates, So when it comes to collecting remotelinks, we only need to compare an issue's `updated` field with its
		`remotelink_updated` field. If `remotelink_updated` is older, then we'll collect remotelinks for this issue and
		set its `remotelink_updated` to `updated` at the end.
	*/
	cursor, err := db.Model(jiraIssue).
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
			data.Options.SourceId,
			data.Options.BoardId,
		).
		Rows()
	if err != nil {
		logger.Error("collect remotelink error:%v", err)
		return err
	}
	defer cursor.Close()

	// iterate all rows
	for cursor.Next() {
		err = db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		issueId := jiraIssue.IssueId
		updated := jiraIssue.Updated

		collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
			RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
				Ctx: taskCtx,
				Params: JiraApiParams{
					SourceId: data.Source.ID,
					BoardId:  data.Options.BoardId,
				},
				Table: RAW_REMOTELINK_TABLE,
			},
			ApiClient:   data.ApiClient,
			UrlTemplate: fmt.Sprintf("api/2/issue/%d/remotelink", issueId),
			ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
				if res.StatusCode == http.StatusNotFound {
					return nil, nil
				}
				blob, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return nil, err
				}
				res.Body.Close()
				return []json.RawMessage{blob}, nil
			},
		})
		if err != nil {
			return err
		}
		err = collector.Execute()
		if err != nil {
			return err
		}
		issue := &models.JiraIssue{SourceId: data.Source.ID, IssueId: issueId}
		err = db.Model(issue).Update("remotelink_updated", updated).Error
		if err != nil {
			return err
		}
	}
	return nil
}
