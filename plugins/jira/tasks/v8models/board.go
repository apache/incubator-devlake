package v8models

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type Board struct {
	ID   uint64 `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (b Board) toToolLayer(sourceId uint64) *models.JiraBoard {
	return &models.JiraBoard{
		SourceId: sourceId,
		BoardId:  b.ID,
		Name:     b.Name,
		Self:     b.Self,
		Type:     b.Type,
	}
}
func (b Board) FromAPI(sourceId uint64, raw json.RawMessage) (interface{}, error) {
	var board Board
	err := json.Unmarshal(raw, &board)
	if err != nil {
		return nil, err
	}
	return board.toToolLayer(sourceId), nil
}
func (b Board) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	return blob, nil
}
