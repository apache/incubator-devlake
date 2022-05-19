package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = CollectChangelogs

const RAW_CHANGELOG_TABLE = "jira_api_changelogs"

func CollectChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		return nil
	}
	db := taskCtx.GetDb()
	// figure out the time range
	since := data.Since

	// filter out issue_ids that needed collection
	tx := db.Table("_tool_jira_board_issues bi").
		Select("bi.issue_id, NOW() AS update_time").
		Joins("LEFT JOIN _tool_jira_issues i ON (bi.connection_id = i.connection_id AND bi.issue_id = i.issue_id)").
		Where("bi.connection_id = ? AND bi.board_id = ? AND (i.changelog_updated IS NULL OR i.changelog_updated < i.updated)", data.Options.ConnectionId, data.Options.BoardId)

	// apply time range if any
	if since != nil {
		tx = tx.Where("i.updated > ?", *since)
	}

	// construct the input iterator
	cursor, err := tx.Rows()
	if err != nil {
		return err
	}
	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	// now, let ApiCollector takes care the rest
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    50,
		Incremental: true,
		Input:       iterator,
		UrlTemplate: "api/3/issue/{{ .Input.IssueId }}/changelog",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Values []json.RawMessage
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Values, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
