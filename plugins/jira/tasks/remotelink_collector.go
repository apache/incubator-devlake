package tasks

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_REMOTELINK_TABLE = "jira_api_remotelinks"

var _ core.SubTaskEntryPoint = CollectRemotelinks

func CollectRemotelinks(taskCtx core.SubTaskContext) error {
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
		Select("_tool_jira_issues.issue_id", "NOW() AS update_time").
		Joins(`LEFT JOIN _tool_jira_board_issues ON (
			_tool_jira_board_issues.connection_id = _tool_jira_issues.connection_id AND
			_tool_jira_board_issues.issue_id = _tool_jira_issues.issue_id
		)`).
		Where(`
			_tool_jira_board_issues.connection_id = ? AND
			_tool_jira_board_issues.board_id = ? AND
			(_tool_jira_issues.remotelink_updated IS NULL OR _tool_jira_issues.remotelink_updated < _tool_jira_issues.updated)
			`,
			data.Options.ConnectionId,
			data.Options.BoardId,
		).
		Rows()
	if err != nil {
		logger.Error("collect remotelink error:%v", err)
		return err
	}
	defer cursor.Close()

	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "api/2/issue/{{ .Input.IssueId }}/remotelink",
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			if res.StatusCode == http.StatusNotFound {
				return nil, nil
			}
			var result []json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	})
	if err != nil {
		return err
	}
	err = collector.Execute()
	if err != nil {
		return err
	}
	return nil
}
