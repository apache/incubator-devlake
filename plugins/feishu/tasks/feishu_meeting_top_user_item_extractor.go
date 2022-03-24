package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/feishu/models"
)

var _ core.SubTaskEntryPoint = ExtractMeetingTopUserItem

func ExtractMeetingTopUserItem(taskCtx core.SubTaskContext) error{
	
	exetractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ApiResName: "top_user_report",
			},
			Table: RAW_MEETING_TOP_USER_ITEM_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error){
			body := &models.FeishuMeetingTopUserItem{}
			err := json.Unmarshal(row.Data, body)
			if err != nil{
				return nil, err
			}
			rawInput := &helper.DatePair{}
			rawErr := json.Unmarshal(row.Input, rawInput)
			if rawErr != nil{
				return nil, rawErr
			}
			results := make([]interface{}, 0)
			results = append(results, &models.FeishuMeetingTopUserItem{
				StartTime: rawInput.PairStartTime.AddDate(0, 0, -1),
				MeetingCount: body.MeetingCount,
				MeetingDuration: body.MeetingDuration,
				Name: body.Name,
				UserType: body.UserType,
			})
			return results, nil
		},	
	})
	if err != nil{
		return err
	}

	return exetractor.Execute()
}

var ExtractMeetingTopUserItemMeta = core.SubTaskMeta{
	Name: "extractMeetingTopUserItem",
	EntryPoint: ExtractMeetingTopUserItem,
	EnabledByDefault: true,
	Description: "Extrat raw top user meeting data into tool layer table feishu_meeting_top_user_item",
}
