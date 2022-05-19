package tasks

import (
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type FeishuApiParams struct {
	ApiResName string `json:"apiResName"`
}

type FeishuOptions struct {
	NumOfDaysToCollect float64  `json:"numOfDaysToCollect"`
	Tasks              []string `json:"tasks,omitempty"`
}

type FeishuTaskData struct {
	Options                  *FeishuOptions
	ApiClient                *helper.ApiAsyncClient
	FeishuMeetingTopUserItem *models.FeishuMeetingTopUserItem
}
