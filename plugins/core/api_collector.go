package core

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"
)

type AsyncApiClient interface {
	GetAsync(path string, queryParams *url.Values, handler func(*http.Response) error) error
}

type Pager struct {
	Page int
	Skip int
	Size int
}

type Iterator interface {
	Fetch() (interface{}, error)
	Close() error
}

type ApiCollectorArgs struct {
	Table         string      `comment:"Raw data table name"`
	UrlTemplate   string      `comment:"GoTemplate for API url"`
	Params        interface{} `comment:"To identify a set of records with same UrlTemplate, i.e. {SourceId, BoardId} for jira entities"`
	Query         *url.Values `comment:"Extra query string when requesting API, like 'Since' option for jira issues collection"`
	ResponseData  interface{}
	PageSize      int
	ApiClient     AsyncApiClient
	SetPage       func(req *http.Request)
	Input         func() Iterator
	OnData        func(res *http.Response, body interface{}) (interface{}, error)
	GetTotalPages func(res *http.Response, body interface{}) (int, error)
}

type ApiCollector struct {
	args        *ApiCollectorArgs
	urlTemplate *template.Template
}

func NewApiCollector(args ApiCollectorArgs) (*ApiCollector, error) {
	if args.Table == "" {
		return nil, fmt.Errorf("Table argument is required")
	}
	tpl, err := template.New(args.Table).Parse(args.UrlTemplate)
	if err != nil {
		return nil, err
	}
	return &ApiCollector{
		args:        &args,
		urlTemplate: tpl,
	}, nil
}

type RawDataTable struct {
	Uuid   string `gorm:"type:uuid;primaryKey"`
	Url    string `gorm:"type:varchar(255)"`
	Params string `gorm:"type:string"`
}

func (collector *ApiCollector) Execute() error {

	return nil
}

var _ SubTask = (*ApiCollector)(nil)
