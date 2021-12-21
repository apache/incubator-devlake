package apimodels

import "github.com/merico-dev/lake/plugins/feishu/models"

type FeishuMeetingTopUserItemResult struct {
	Code int64 `json:"code"`
	Data struct {
		TopUserReport []models.FeishuMeetingTopUserItem `json:"top_user_report"`
	} `json:"data"`
	Msg string `json:"msg"`
}
