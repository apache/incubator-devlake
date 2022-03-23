package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
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
	incremental := false
	if since == nil {
		// user didn't specify a time range to sync, try load from database
		var latestUpdated models.JiraChangelog
		err := db.Where("source_id = ?", data.Source.ID).Order("created DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return fmt.Errorf("failed to get latest jira changelog record: %w", err)
		}
		if latestUpdated.ChangelogId > 0 {
			since = &latestUpdated.Created
			incremental = true
		}
	}

	// filter out issue_ids that needed collection
	tx := db.Table("jira_board_issues bi").
		Select("bi.issue_id, NOW() AS update_time").
		Joins("LEFT JOIN jira_issues i ON (bi.source_id = i.source_id AND bi.issue_id = i.issue_id)").
		Where("bi.source_id = ? AND bi.board_id = ? AND (i.changelog_updated IS NULL OR i.changelog_updated < i.updated)", data.Options.SourceId, data.Options.BoardId)

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
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    50,
		Incremental: incremental,
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
