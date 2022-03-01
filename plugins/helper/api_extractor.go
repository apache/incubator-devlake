package helper

import "github.com/merico-dev/lake/plugins/core"

type RawDataExtractor func(row interface{}, params interface{}) (interface{}, error)

type ApiExtractorArgs struct {
	Table      string `comment:"Raw data table name"`
	Url        string `comment:"Raw data table name"`
	Params     interface{}
	RowData    interface{}
	Extractors []RawDataExtractor
	BatchSize  int
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
