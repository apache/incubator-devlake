package helper

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/merico-dev/lake/plugins/core"
)

var _ core.SubTask = (*ApiSyncCollector)(nil)

// SyncResponseHandler returns a bool value indicates if the iteration should continue
type SyncResponseHandler func(res *http.Response) (shouldContinue bool, err error)
type ApiSyncCollectorArgs struct {
	RawDataSubTaskArgs
	UrlTemplate     string                                 `comment:"GoTemplate for API url"`
	Query           func(pager *Pager) (url.Values, error) `comment:"Extra query string when requesting API, like 'Since' option for jira issues collection"`
	Header          func(pager *Pager) (http.Header, error)
	PageSize        int
	Incremental     bool `comment:"Indicate this is a incremental collection, so the existing data won't get flushed"`
	ApiClient       core.SyncApiClient
	Input           Iterator
	InputRateLimit  int
	responseHandler SyncResponseHandler
}

type ApiSyncCollector struct {
	*RawDataSubTask
	args        *ApiSyncCollectorArgs
	urlTemplate *template.Template
	nextTime    time.Time
	waitTime    time.Duration
}

func NewApiSyncCollector(args ApiSyncCollectorArgs) (*ApiSyncCollector, error) {
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
	return &ApiSyncCollector{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
		urlTemplate:    tpl,
		nextTime:       time.Now(),
		waitTime:       time.Second / time.Duration(args.InputRateLimit),
	}, nil
}

// Start collection
func (collector *ApiSyncCollector) Execute() error {
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
			err = collector.fetchPages(input)
			if err != nil {
				break
			}
		}

	} else {
		// or we just did it once
		err = collector.fetchPages(nil)
	}

	logger.Info("end api collection")
	return err
}

func (collector *ApiSyncCollector) generateUrl(pager *Pager, input interface{}) (string, error) {
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

func (collector *ApiSyncCollector) fetchPages(input interface{}) error {
	var err error
	var shouldContinue bool
	for page := 1; shouldContinue; page++ {
		pager := &Pager{
			Page: page,
			Size: collector.args.PageSize,
			Skip: collector.args.PageSize * (page - 1),
		}
		select {
		case <-collector.args.Ctx.GetContext().Done():
			return collector.args.Ctx.GetContext().Err()
		default:
		}
		// rate limit
		<-time.NewTimer(time.Since(collector.nextTime)).C

		collector.nextTime = collector.nextTime.Add(collector.waitTime)
		shouldContinue, err = collector.fetchSync(pager, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (collector *ApiSyncCollector) fetchSync(pager *Pager, input interface{}) (bool, error) {
	if pager == nil {
		pager = &Pager{
			Page: 1,
			Size: 100,
			Skip: 0,
		}
	}
	apiUrl, err := collector.generateUrl(pager, input)
	if err != nil {
		return false, err
	}
	var apiQuery url.Values
	if collector.args.Query != nil {
		apiQuery, err = collector.args.Query(pager)
		if err != nil {
			return false, err
		}
	}
	apiHeader := (http.Header)(nil)
	if collector.args.Header != nil {
		apiHeader, err = collector.args.Header(pager)
		if err != nil {
			return false, err
		}
	}

	return collector.args.ApiClient.GetSync(apiUrl, apiQuery, apiHeader, collector.args.responseHandler)
}
