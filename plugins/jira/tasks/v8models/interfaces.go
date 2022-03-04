package v8models

import (
	"encoding/json"
)

type Transformer interface {
	FromAPI(sourceId uint64, raw json.RawMessage) (interface{}, error)
	ExtractRawMessage(blob []byte) (json.RawMessage, error)
}
type TransformerWithIssueId interface {
	FromAPI(sourceId, issueId uint64, raw json.RawMessage) (interface{}, error)
	ExtractRawMessage(blob []byte) (json.RawMessage, error)
}
