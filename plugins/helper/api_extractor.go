package helper

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
)

// Accept raw json body and params, return list of entities that need to be stored
type RawDataExtractor func(body json.RawMessage, params json.RawMessage) ([]interface{}, error)

type ApiExtractorArgs struct {
	Ctx       core.TaskContext
	Table     string `comment:"Raw data table name"`
	Url       string `comment:"Raw data table name"`
	Params    interface{}
	RowData   interface{}
	Extract   RawDataExtractor
	BatchSize int
}

type ApiExtractor struct {
	args *ApiExtractorArgs
}

func NewApiExtractor(args ApiExtractorArgs) (*ApiExtractor, error) {
	return &ApiExtractor{args: &args}, nil
}

func (extractor *ApiExtractor) Execute() error {
	return nil
}

var _ core.SubTask = (*ApiExtractor)(nil)
