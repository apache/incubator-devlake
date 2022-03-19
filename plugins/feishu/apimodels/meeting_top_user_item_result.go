package apimodels

import "encoding/json"

type FeishuMeetingTopUserItemResult struct {
	Code int64 `json:"code"`
	Data struct {
		TopUserReport []json.RawMessage `json:"top_user_report"`
	} `json:"data"`
	Msg string `json:"msg"`
}
