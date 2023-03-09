package tasks

import (
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TrelloOptions struct {
	ConnectionId         uint64 `json:"connectionId"`
	BoardId              string `json:"boardId"`
	TransformationRuleId uint64
}

type TrelloTaskData struct {
	Options   *TrelloOptions
	ApiClient *api.ApiAsyncClient
}

type TrelloApiParams struct {
	ConnectionId uint64
	BoardId      string
}
