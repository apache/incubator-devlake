package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
	"gorm.io/datatypes"
)

type Pager struct {
	Page int
	Skip int
	Size int
}

type UrlData struct {
	Pager  *Pager
	Params interface{}
	Input  interface{}
}

type AsyncResponseHandler func(res *http.Response) error

type ApiCollectorArgs struct {
	RawDataSubTaskArgs
	UrlTemplate    string                                 `comment:"GoTemplate for API url"`
	Query          func(pager *Pager) (url.Values, error) `comment:"Extra query string when requesting API, like 'Since' option for jira issues collection"`
	Header         func(pager *Pager) (http.Header, error)
	PageSize       int
	Incremental    bool `comment:"Indicate this is a incremental collection, so the existing data won't get flushed"`
	ApiClient      core.AsyncApiClient
	Input          Iterator
	InputRateLimit int
	GetTotalPages  func(res *http.Response, args *ApiCollectorArgs) (int, error)
}

type ApiCollector struct {
	*RawDataSubTask
	args        *ApiCollectorArgs
	urlTemplate *template.Template
}

// NewApiCollector allocates a new ApiCollector  with the given args.
// ApiCollector can help you collecting data from some api with ease, pass in a AsyncApiClient and tell it which part
// of response you want to save, ApiCollector will collect them from remote server and store them into database.
func NewApiCollector(args ApiCollectorArgs) (*ApiCollector, error) {
	// process args
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	// TODO: check if args.Table is valid
	if args.UrlTemplate == "" {
		return nil, fmt.Errorf("UrlTemplate is required")
	}
	tpl, err := template.New(args.Table).Parse(args.UrlTemplate)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile UrlTemplate: %w", err)
	}
	if args.ApiClient == nil {
		return nil, fmt.Errorf("ApiClient is required")
	}
	if args.InputRateLimit == 0 {
		args.InputRateLimit = 50
	}
	return &ApiCollector{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
		urlTemplate:    tpl,
	}, nil
}

// Start collection
func (collector *ApiCollector) Execute() error {
	logger := collector.args.Ctx.GetLogger()
	logger.Info("start api collection")

	// make sure table is created
	db := collector.args.Ctx.GetDb()
	err := db.Table(collector.table).AutoMigrate(&RawData{})
	if err != nil {
		return err
	}

	// flush data if not incremental collection
	if !collector.args.Incremental {
		err = db.Table(collector.table).Delete(&RawData{}, "params = ?", collector.params).Error
		if err != nil {
			return err
		}
	}

	if collector.args.Input != nil {
		// if Input was given, we iterate through it and exec multiple times
		// create a parent scheduler, note that the rate limit of this scheduler is different than
		// api rate limit
		scheduler, err := utils.NewWorkerScheduler(
			collector.args.InputRateLimit*6/5, // increase by 20 percent
			collector.args.InputRateLimit,
			collector.args.Ctx.GetContext(),
		)
		if err != nil {
			return err
		}
		defer scheduler.Release()

		collector.args.Ctx.SetProgress(0, -1)
		// load all rows from iterator, and exec them in parallel
		// TODO: this loads all records into memory, we need lazy-load
		iterator := collector.args.Input
		defer iterator.Close()
		for iterator.HasNext() {
			input, err := iterator.Fetch()
			if err != nil {
				return err
			}
			err = collector.exec(input)
			if err != nil {
				break
			}
		}

		scheduler.WaitUntilFinish()
	} else {
		// or we just did it once
		err = collector.exec(nil)
	}

	collector.args.ApiClient.WaitAsync()
	logger.Info("end api collection")
	return err
}

func (collector *ApiCollector) exec(input interface{}) error {
	if collector.args.PageSize > 0 {
		// collect multiple pages
		return collector.fetchPagesAsync(input)
	}
	// collect detail of a record
	return collector.fetchAsync(nil, input, collector.handleResponse)
}

func (collector *ApiCollector) generateUrl(pager *Pager, input interface{}) (string, error) {
	var buf bytes.Buffer
	err := collector.urlTemplate.Execute(&buf, &UrlData{
		Pager:  pager,
		Params: collector.args.Params,
		Input:  input,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (collector *ApiCollector) fetchPagesAsync(input interface{}) error {
	var err error
	if collector.args.GetTotalPages != nil {
		/* when total pages is available from api*/
		// fetch the very first page
		err = collector.fetchAsync(nil, input, func(res *http.Response) error {
			// gather total pages
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			res.Body.Close()
			res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			totalPages, err := collector.args.GetTotalPages(res, collector.args)
			if err != nil {
				return err
			}
			// save response body of first page
			res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			err = collector.handleResponse(res)
			if err != nil {
				return err
			}
			if collector.args.Input == nil {
				collector.args.Ctx.SetProgress(1, totalPages)
			}
			// fetch other pages in parallel
			for page := 2; page <= totalPages; page++ {
				err = collector.fetchAsync(&Pager{
					Page: page,
					Size: collector.args.PageSize,
					Skip: collector.args.PageSize * (page - 1),
				}, input, func(res *http.Response) error {
					err := collector.handleResponse(res)
					if err != nil {
						return err
					}
					if collector.args.Input == nil {
						collector.args.Ctx.IncProgress(1)
					}
					return nil
				})
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
	if err != nil {
		return err
	}
	if collector.args.Input != nil {
		collector.args.Ctx.IncProgress(1)
	}
	return nil
}

func (collector *ApiCollector) handleResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	db := collector.args.Ctx.GetDb()
	return db.Table(collector.table).Create(&RawData{
		Data:   datatypes.JSON(body),
		Params: collector.params,
		Url:    res.Request.URL.String(),
	}).Error
}

func (collector *ApiCollector) fetchAsync(pager *Pager, input interface{}, handler AsyncResponseHandler) error {
	if pager == nil {
		pager = &Pager{
			Page: 1,
			Size: 100,
			Skip: 0,
		}
	}
	apiUrl, err := collector.generateUrl(pager, input)
	if err != nil {
		return err
	}
	apiQuery := (url.Values)(nil)
	if collector.args.Query != nil {
		apiQuery, err = collector.args.Query(pager)
		if err != nil {
			return err
		}
	}
	apiHeader := (http.Header)(nil)
	if collector.args.Header != nil {
		apiHeader, err = collector.args.Header(pager)
		if err != nil {
			return err
		}
	}
	return collector.args.ApiClient.GetAsync(apiUrl, apiQuery, apiHeader, handler)
}

var _ core.SubTask = (*ApiCollector)(nil)
