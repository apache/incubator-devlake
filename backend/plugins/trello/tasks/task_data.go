package tasks

import (
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"time"
)

type TrelloOptions struct {
	ConnectionId         uint64 `json:"connectionId"`
	BoardId              string `json:"boardId"`
	TimeAfter            string
	TransformationRuleId uint64
}

type TrelloTaskData struct {
	Options   *TrelloOptions
	ApiClient *api.ApiAsyncClient
	TimeAfter *time.Time
}

type TrelloApiParams struct {
	ConnectionId uint64
	BoardId      string
}
