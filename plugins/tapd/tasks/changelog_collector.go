package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_CHANGELOG_TABLE = "tapd_api_changelogs"

var _ core.SubTaskEntryPoint = CollectChangelogs

type Type struct {
	Type string
}

func CollectChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	logger := taskCtx.GetLogger()
	logger.Info("collect storyChangelogs")

	cursor, err := db.Raw("select 'task' as 'type' UNION select 'story' as 'type' UNION select 'bug' as 'type'").Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(Type{}))
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		Input:       iterator,
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "{{ .Input.Type }}_changes",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			//query.Set("order", "created,asc")
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error("collect story changelog error:", err)
		return err
	}
	return collector.Execute()
}

var CollectChangelogMeta = core.SubTaskMeta{
	Name:        "collectStoryChangelogs",
	EntryPoint:  CollectChangelogs,
	Required:    true,
	Description: "collect Tapd storyChangelogs",
}
