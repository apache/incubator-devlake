package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

const RAW_WORKLOGS_TABLE = "jira_api_worklogs"

func CollectWorklogs(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JiraTaskData)
	since := data.Since
	incremental := false

	if since == nil {
		var latestUpdated models.JiraWorklog
		err := db.Where("connection_id = ?", data.Connection.ID).Order("updated DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return fmt.Errorf("failed to get latest jira issue worklog record: %w", err)
		}
		if latestUpdated.IssueId > 0 {
			since = &latestUpdated.Updated
			incremental = true
		}
	}

	logger := taskCtx.GetLogger()
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	tx := db.Model(&models.JiraIssue{}).
		Joins("left join _tool_jira_board_issues on _tool_jira_issues.issue_id = _tool_jira_board_issues.issue_id").
		Select("_tool_jira_board_issues.issue_id").Where("_tool_jira_board_issues.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId)

	if since != nil {
		tx = tx.Where("_tool_jira_issues.updated > ?", since)
	}
	cursor, err := tx.Rows()
	if err != nil {
		return err
	}
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
			Table: RAW_WORKLOGS_TABLE,
		},
		Input:         iterator,
		ApiClient:     data.ApiClient,
		UrlTemplate:   "api/2/issue/{{ .Input.IssueId }}/worklog",
		PageSize:      50,
		Incremental:   incremental,
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Worklogs []json.RawMessage `json:"worklogs"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Worklogs, nil
		},
	})
	if err != nil {
		logger.Error("collect board error:", err)
		return err
	}

	return collector.Execute()
}
