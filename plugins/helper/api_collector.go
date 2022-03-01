package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"text/template"

	"github.com/google/uuid"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/datatypes"
)

type Pager struct {
	Page int
	Skip int
	Size int
}

type Iterator interface {
	Fetch() (interface{}, error)
	Close() error
}

type UrlData struct {
	Pager  *Pager
	Params interface{}
}

type AsyncResponseHandler func(res *http.Response, body interface{}) error

type ApiCollectorArgs struct {
	Ctx           core.TaskContext
	Table         string                                  `comment:"Raw data table name"`
	UrlTemplate   string                                  `comment:"GoTemplate for API url"`
	Params        interface{}                             `comment:"To identify a set of records with same UrlTemplate, i.e. {SourceId, BoardId} for jira entities"`
	Query         func(pager *Pager) (*url.Values, error) `comment:"Extra query string when requesting API, like 'Since' option for jira issues collection"`
	Header        func(pager *Pager) (*url.Values, error)
	PageSize      int
	ApiClient     core.AsyncApiClient
	Input         func() Iterator
	BodyType      reflect.Type
	OnData        func(res *http.Response, body interface{}) (interface{}, error)
	GetTotalPages func(res *http.Response, body interface{}) (int, error)
}

type ApiCollector struct {
	args        *ApiCollectorArgs
	urlTemplate *template.Template
	table       string
	params      string
}

// NewApiCollector allocates a new ApiCollector  with the given args.
// ApiCollector can help you collecting data from some api with ease, pass in a AsyncApiClient and tell it which part
// of response you want to save, ApiCollector will collect them from remote server and store them into database.
func NewApiCollector(args ApiCollectorArgs) (*ApiCollector, error) {
	// process args
	if args.Ctx == nil {
		return nil, fmt.Errorf("Ctx is required")
	}
	if args.Table == "" {
		return nil, fmt.Errorf("Table is required")
	}
	// TODO: check if args.Table is valid
	if args.UrlTemplate == "" {
		return nil, fmt.Errorf("UrlTemplate is required")
	}
	tpl, err := template.New(args.Table).Parse(args.UrlTemplate)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile UrlTemplate: %w", err)
	}
	params := ""
	if args.Params == nil {
		args.Ctx.GetLogger().Warn("Missing `Params` for collecting %s to %s", args.UrlTemplate, args.Table)
	} else {
		// TODO: sort it to make it consisitence
		paramsBytes, err := json.Marshal(args.Params)
		if err != nil {
			return nil, err
		}
		params = string(paramsBytes)
	}
	if args.ApiClient == nil {
		return nil, fmt.Errorf("ApiClient is required")
	}
	if args.BodyType == nil {
		return nil, fmt.Errorf("BodyType is required")
	}
	if args.OnData == nil {
		return nil, fmt.Errorf("OnData handler is required")
	}
	return &ApiCollector{
		args:        &args,
		urlTemplate: tpl,
		table:       fmt.Sprintf("_raw_%s", args.Table),
		params:      params,
	}, nil
}

// Start collection
func (collector *ApiCollector) Execute() error {
	// make sure table is created
	db := collector.args.Ctx.GetDb()
	for _, creationSql := range GetRawTableCreationSqls(collector.table) {
		err := db.Exec(creationSql).Error
		if err != nil {
			return err
		}
	}

	// flush data, TODO: incremental data collection
	err := db.Exec(GetRawTableDeletionSql(collector.table), collector.params).Error
	if err != nil {
		return err
	}
	if collector.args.PageSize > 0 {
		// collect multiple pages
		err = collector.fetchPagesAsync()
	} else {
		// collect detail of a record
		err = collector.fetchAsync(nil, collector.handleResponse)
	}
	if err != nil {
		return err
	}
	collector.args.ApiClient.WaitAsync()
	return nil
}

func (collector *ApiCollector) generateUrl(pager *Pager) (string, error) {
	var buf bytes.Buffer
	err := collector.urlTemplate.Execute(&buf, &UrlData{
		Pager:  pager,
		Params: collector.args.Params,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (collector *ApiCollector) newResponseBody() interface{} {
	return reflect.New(collector.args.BodyType.Elem()).Interface()
}

func (collector *ApiCollector) fetchPagesAsync() error {
	if collector.args.GetTotalPages != nil {
		/* when total pages is available from api*/
		// fetch the very first page
		return collector.fetchAsync(nil, func(res *http.Response, body interface{}) error {
			// gather total pages
			totalPages, err := collector.args.GetTotalPages(res, body)
			if err != nil {
				return err
			}
			// handle response body of first page
			err = collector.handleResponse(res, body) // THIS WONT WORK: because res.Body already got read in GetTotalPages
			if err != nil {
				return err
			}
			// fetch other pages in parallel
			for page := 2; page <= totalPages; page++ {
				err = collector.fetchAsync(&Pager{
					Page: page,
					Size: collector.args.PageSize,
					Skip: collector.args.PageSize * (page - 1),
				}, collector.handleResponse)
				if err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		/* when total pages is available from api*/
		// fetch page by page in sequential?
		// use step currency technique? fetch like 10 pages at once, if all went well, fetch next 10 pages?
		panic("not implmented")
	}
}

func (collector *ApiCollector) handleResponse(res *http.Response, body interface{}) error {
	//
	data, err := collector.args.OnData(res, body)
	if err != nil {
		return err
	}
	db := collector.args.Ctx.GetDb()
	switch d := data.(type) {
	case []json.RawMessage:
		rows := make([]*RawData, len(d))
		for i, item := range d {
			rows[i] = &RawData{
				Uuid:   uuid.New(),
				Data:   datatypes.JSON(item),
				Params: collector.params,
			}
		}
		err = db.Table(collector.table).CreateInBatches(rows, len(rows)).Error
	case json.RawMessage:
		err = db.Table(collector.table).Create(&RawData{
			Uuid:   uuid.New(),
			Data:   datatypes.JSON(d),
			Params: collector.params,
		}).Error
	default:
		err = fmt.Errorf("unexpected type returned")
	}
	return err
}

func (collector *ApiCollector) fetchAsync(pager *Pager, handler AsyncResponseHandler) error {
	if pager == nil {
		pager = &Pager{
			Page: 1,
			Size: 100,
			Skip: 0,
		}
	}
	apiUrl, err := collector.generateUrl(pager)
	if err != nil {
		return err
	}
	apiQuery := (*url.Values)(nil)
	if collector.args.Query != nil {
		apiQuery, err = collector.args.Query(pager)
		if err != nil {
			return err
		}
	}
	apiHeader := (*url.Values)(nil)
	if collector.args.Header != nil {
		apiHeader, err = collector.args.Header(pager)
		if err != nil {
			return err
		}
	}
	return collector.args.ApiClient.GetAsync(apiUrl, apiQuery, apiHeader, func(res *http.Response) error {
		body := collector.newResponseBody()
		err := core.UnmarshalResponse(res, body)
		if err != nil {
			return err
		}
		return handler(res, body)
	})
}

var _ core.SubTask = (*ApiCollector)(nil)
