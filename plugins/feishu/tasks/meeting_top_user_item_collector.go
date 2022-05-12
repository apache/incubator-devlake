package tasks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/feishu/apimodels"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_MEETING_TOP_USER_ITEM_TABLE = "feishu_meeting_top_user_item"

var _ core.SubTaskEntryPoint = CollectMeetingTopUserItem

func CollectMeetingTopUserItem(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*FeishuTaskData)
	pageSize := 100
	NumOfDaysToCollectInt := int(data.Options.NumOfDaysToCollect)
	iterator, err := helper.NewDateIterator(NumOfDaysToCollectInt)
	if err != nil {
		return err
	}
	incremental := false

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ApiResName: "top_user_report",
			},
			Table: RAW_MEETING_TOP_USER_ITEM_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: incremental,
		Input:       iterator,
		UrlTemplate: "/reports/get_top_user",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			input := reqData.Input.(*helper.DatePair)
			query.Set("start_time", strconv.FormatInt(input.PairStartTime.Unix(), 10))
			query.Set("end_time", strconv.FormatInt(input.PairEndTime.Unix(), 10))
			query.Set("limit", strconv.Itoa(pageSize))
			query.Set("order_by", "2")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body := &apimodels.FeishuMeetingTopUserItemResult{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			return body.Data.TopUserReport, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectMeetingTopUserItemMeta = core.SubTaskMeta{
	Name: "collectMeetingTopUserItem",
	EntryPoint: CollectMeetingTopUserItem,
	EnabledByDefault: true,
	Description: "Collect top user meeting data from Feishu api",
}

