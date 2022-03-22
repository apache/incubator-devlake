package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/jira/models"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_WORKLOGS_TABLE = "jira_api_worklogs"

// make sure table _raw_jira_api_worklogs exists
func CollectWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	var boardIssue models.JiraBoardIssue
	err := db.First(&boardIssue, "source_id = ? AND board_id = ?", sourceId, boardId).Error
	if err != nil {
		return err
	}
	logger.Info("collect worklog, board_id:%d, issue_id:%d", data.Options.BoardId, boardIssue.IssueId)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: fmt.Sprintf("api/2/issue/%d/worklog", boardIssue.IssueId),
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			return nil, nil
		},
	})
	if err != nil {
		logger.Error("collect board error:", err)
		return err
	}

	return collector.Execute()
}

func collectWorklogs(taskCtx core.SubTaskContext, issueId uint64) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect worklog, board_id:%d, issue_id:%d", data.Options.BoardId, issueId)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: fmt.Sprintf("api/2/issue/%d/worklog", issueId),
		PageSize:    50,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var result []json.RawMessage
			err := core.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	})
	if err != nil {
		logger.Error("collect board error:", err)
		return err
	}

	return collector.Execute()
}
