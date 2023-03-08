package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
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

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*TrelloOptions, errors.Error) {
	var op TrelloOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to decode trello options")
	}

	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("trello connectionId is invalid")
	}
	if op.BoardId == "" {
		return nil, errors.BadInput.New(fmt.Sprintf("trello boardId is invalid"))
	}
	return &op, nil
}
