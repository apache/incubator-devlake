/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/merico-dev/graphql"
)

// CursorPager contains pagination information for a graphql request
type CursorPager struct {
	SkipCursor *string
	Size       int
}

// GraphqlRequestData is the input of `UrlTemplate` `BuildQuery` and `Header`, so we can generate them dynamically
type GraphqlRequestData struct {
	Pager     *CursorPager
	Params    interface{}
	Input     interface{}
	InputJSON []byte
}

// GraphqlQueryPageInfo contains the pagination data
type GraphqlQueryPageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

// DateTime is the type of time in Graphql
// graphql lib can only read this name...
type DateTime struct{ time.Time }

// GraphqlAsyncResponseHandler callback function to handle the Response asynchronously
type GraphqlAsyncResponseHandler func(res *http.Response) error

// GraphqlCollectorArgs arguments needed by GraphqlCollector
type GraphqlCollectorArgs struct {
	RawDataSubTaskArgs
	// BuildQuery would be sent out as part of the request URL
	BuildQuery func(reqData *GraphqlRequestData) (query interface{}, variables map[string]interface{}, err error)
	// PageSize tells ApiCollector the page size
	PageSize int
	// GraphqlClient is a asynchronize api request client with qps
	GraphqlClient *GraphqlAsyncClient
	// Input helps us collect data based on previous collected data, like collecting changelogs based on jira
	// issue ids
	Input Iterator
	// how many times fetched from input, default 1 means only fetch once
	// NOTICE: InputStep=1 will fill value as item and InputStep>1 will fill value as []item
	InputStep int
	// Incremental indicate if this is a incremental collection, the existing data won't get deleted if it was true
	Incremental bool `comment:"indicate if this collection is incremental update"`
	// GetPageInfo is to tell `GraphqlCollector` is page information
	GetPageInfo func(query interface{}, args *GraphqlCollectorArgs) (*GraphqlQueryPageInfo, error)
	BatchSize   int
	// one of ResponseParser and ResponseParserEvenWhenDataErrors is required to parse response
	ResponseParser               func(query interface{}, variables map[string]interface{}) ([]interface{}, error)
	ResponseParserWithDataErrors func(query interface{}, variables map[string]interface{}, dataErrors []graphql.DataError) ([]interface{}, error)
}

// GraphqlCollector help you collect data from Graphql services
type GraphqlCollector struct {
	*RawDataSubTask
	args         *GraphqlCollectorArgs
	workerErrors []error
}

// ErrFinishCollect is an error which will finish this collector
var ErrFinishCollect = errors.Default.New("finish collect")

// NewGraphqlCollector allocates a new GraphqlCollector with the given args.
// GraphqlCollector can help us collect data from api with ease, pass in a AsyncGraphqlClient and tell it which part
// of response we want to save, GraphqlCollector will collect them from remote server and store them into database.
func NewGraphqlCollector(args GraphqlCollectorArgs) (*GraphqlCollector, errors.Error) {
	// process args
	rawDataSubTask, err := NewRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	if args.GraphqlClient == nil {
		return nil, errors.Default.New("ApiClient is required")
	}
	if args.ResponseParser == nil && args.ResponseParserWithDataErrors == nil {
		return nil, errors.Default.New("one of ResponseParser and ResponseParserWithDataErrors is required")
	}
	apiCollector := &GraphqlCollector{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	if args.InputStep == 0 {
		args.InputStep = 1
	}
	return apiCollector, nil
}

// Execute api collection
func (collector *GraphqlCollector) Execute() errors.Error {
	logger := collector.args.Ctx.GetLogger()
	logger.Info("start graphql collection")

	// make sure table is created
	db := collector.args.Ctx.GetDal()
	err := db.AutoMigrate(&RawData{}, dal.From(collector.table))
	if err != nil {
		return errors.Default.Wrap(err, "error running auto-migrate")
	}

	divider := NewBatchSaveDivider(collector.args.Ctx, collector.args.BatchSize, collector.table, collector.params)

	isIncremental := collector.args.Incremental
	syncPolicy := collector.args.Ctx.TaskContext().SyncPolicy()
	if syncPolicy != nil && syncPolicy.FullSync {
		isIncremental = false
	}
	// flush data if not incremental collection
	if isIncremental {
		// re extract data for new scope config
		err = collector.ExtractExistRawData(divider)
		if err != nil {
			collector.checkError(err)
		}
	} else {
		err = db.Delete(&RawData{}, dal.From(collector.table), dal.Where("params = ?", collector.params))
		if err != nil {
			return errors.Default.Wrap(err, "error deleting data from collector")
		}
	}

	collector.args.Ctx.SetProgress(0, -1)
	if collector.args.Input != nil {
		iterator := collector.args.Input
		defer iterator.Close()
		// the comment about difference is written at GraphqlCollectorArgs.InputStep
		if collector.args.InputStep == 1 {
			for iterator.HasNext() && !collector.HasError() {
				input, err := iterator.Fetch()
				if err != nil {
					collector.checkError(err)
					break
				}
				collector.exec(divider, input)
			}
		} else {
			for !collector.HasError() {
				var inputs []interface{}
				for i := 0; i < collector.args.InputStep && iterator.HasNext(); i++ {
					input, err := iterator.Fetch()
					if err != nil {
						collector.checkError(err)
						break
					}
					inputs = append(inputs, input)
				}
				if inputs == nil {
					break
				}
				collector.exec(divider, inputs)
			}
		}
	} else {
		// or we just did it once
		collector.exec(divider, nil)
	}

	logger.Debug("wait for all async api to finished")
	collector.args.GraphqlClient.Wait()

	if collector.HasError() {
		err = errors.Default.Combine(collector.workerErrors)
		logger.Error(err, "ended Graphql collector execution with error")
		logger.Error(collector.workerErrors[0], "the first error of them")
		return err
	} else {
		logger.Info("ended api collection without error")
	}

	err = divider.Close()
	return err
}

func (collector *GraphqlCollector) exec(divider *BatchSaveDivider, input interface{}) {
	inputJson, err := json.Marshal(input)
	if err != nil {
		collector.checkError(errors.Default.Wrap(err, `input can not be marshal to json`))
	}
	reqData := new(GraphqlRequestData)
	reqData.Input = input
	reqData.InputJSON = inputJson
	reqData.Pager = &CursorPager{
		SkipCursor: nil,
		Size:       collector.args.PageSize,
	}
	if collector.args.GetPageInfo != nil {
		collector.fetchOneByOne(divider, reqData)
	} else {
		collector.fetchAsync(divider, reqData, nil)
	}
}

// fetchOneByOne fetches data of all pages for APIs that return paging information
func (collector *GraphqlCollector) fetchOneByOne(divider *BatchSaveDivider, reqData *GraphqlRequestData) {
	// fetch first page
	var fetchNextPage func(query interface{}) errors.Error
	fetchNextPage = func(query interface{}) errors.Error {
		pageInfo, err := collector.args.GetPageInfo(query, collector.args)
		if err != nil {
			return errors.Default.Wrap(err, "fetchPagesDetermined get totalPages failed")
		}
		if pageInfo == nil {
			return errors.Default.New("fetchPagesDetermined got pageInfo is nil")
		}
		if pageInfo.HasNextPage {
			collector.args.GraphqlClient.NextTick(func() errors.Error {
				reqDataTemp := &GraphqlRequestData{
					Pager: &CursorPager{
						SkipCursor: &pageInfo.EndCursor,
						Size:       collector.args.PageSize,
					},
					Input:     reqData.Input,
					InputJSON: reqData.InputJSON,
				}
				collector.fetchAsync(divider, reqDataTemp, fetchNextPage)
				return nil
			}, collector.checkError)
		}
		return nil
	}
	collector.fetchAsync(divider, reqData, fetchNextPage)
}

// BatchSaveWithOrigin save the results and fill raw data origin for them
func (collector *GraphqlCollector) BatchSaveWithOrigin(divider *BatchSaveDivider, results []interface{}, row *RawData) errors.Error {
	// batch save divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	for _, result := range results {
		// get the batch operator for the specific type
		batch, err := divider.ForType(reflect.TypeOf(result))
		if err != nil {
			return err
		}
		// set raw data origin field
		origin := reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN)
		if origin.IsValid() && origin.IsZero() {
			origin.Set(reflect.ValueOf(common.RawDataOrigin{
				RawDataTable:  collector.table,
				RawDataId:     row.ID,
				RawDataParams: row.Params,
			}))
		}
		// records get saved into db when slots were max outed
		err = batch.Add(result)
		if err != nil {
			return errors.Default.Wrap(err, "error adding result to batch")
		}
	}
	return nil
}

// ExtractExistRawData will extract data from existing data from raw layer if increment
func (collector *GraphqlCollector) ExtractExistRawData(divider *BatchSaveDivider) errors.Error {
	// load data from database
	db := collector.args.Ctx.GetDal()
	logger := collector.args.Ctx.GetLogger()

	clauses := []dal.Clause{
		dal.From(collector.table),
		dal.Where("params = ?", collector.params),
		dal.Orderby("id ASC"),
	}

	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of clauses")
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error running DB query")
	}
	logger.Info("get data from %s where params=%s and got %d", collector.table, collector.params, count)
	defer cursor.Close()

	// progress
	collector.args.Ctx.SetProgress(0, -1)
	ctx := collector.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		// get the type of query and variables. For each iteration, the query should be a different object
		query, variables, _ := collector.args.BuildQuery(nil)
		row := &RawData{}
		err = db.Fetch(cursor, row)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching row")
		}

		err = errors.Convert(json.Unmarshal(row.Data, &query))
		if err != nil {
			return errors.Default.Wrap(err, `graphql collector unmarshal query failed`)
		}
		err = errors.Convert(json.Unmarshal(row.Input, &variables))
		if err != nil {
			return errors.Default.Wrap(err, `variables in graphql query can not unmarshal from json`)
		}

		var results []interface{}
		if collector.args.ResponseParserWithDataErrors != nil {
			results, err = errors.Convert01(collector.args.ResponseParserWithDataErrors(query, variables, nil))
		} else {
			results, err = errors.Convert01(collector.args.ResponseParser(query, variables))
		}
		if err != nil {
			if errors.Is(err, ErrFinishCollect) {
				logger.Info("existing data parser return ErrFinishCollect, but skip. rawId: #%d", row.ID)
			} else {
				return errors.Default.Wrap(err, "error calling plugin Extract implementation")
			}
		}
		err = collector.BatchSaveWithOrigin(divider, results, row)
		if err != nil {
			return err
		}

		collector.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

func (collector *GraphqlCollector) fetchAsync(divider *BatchSaveDivider, reqData *GraphqlRequestData, handler func(query interface{}) errors.Error) {
	if reqData.Pager == nil {
		reqData.Pager = &CursorPager{
			SkipCursor: nil,
			Size:       collector.args.PageSize,
		}
	}
	query, variables, err := collector.args.BuildQuery(reqData)
	if err != nil {
		collector.checkError(errors.Default.Wrap(err, `graphql collector BuildQuery failed`))
		return
	}

	logger := collector.args.Ctx.GetLogger()
	dataErrors, err := collector.args.GraphqlClient.Query(query, variables)
	if err != nil {
		if err == context.Canceled {
			// direct error message for error combine
			collector.checkError(err)
		} else {
			collector.checkError(errors.Default.Wrap(err, `graphql query failed`))
		}
		return
	}
	if len(dataErrors) > 0 {
		if collector.args.ResponseParserWithDataErrors == nil {
			for _, dataError := range dataErrors {
				collector.checkError(errors.Default.Wrap(dataError, `graphql query got error`))
			}
			return
		}
		// else: error will deal by ResponseParserWithDataErrors
	}
	defer logger.Debug("fetchAsync >>> done for %v %v", query, variables)

	paramsBytes, err := json.Marshal(query)
	if err != nil {
		collector.checkError(errors.Default.Wrap(err, `graphql collector marshal query failed`))
		return
	}
	db := collector.args.Ctx.GetDal()
	queryStr, _ := graphql.ConstructQuery(query, variables)
	variablesJson, err := json.Marshal(variables)
	if err != nil {
		collector.checkError(errors.Default.Wrap(err, `variables in graphql query can not marshal to json`))
		return
	}
	row := &RawData{
		Params: collector.params,
		Data:   paramsBytes,
		Url:    queryStr,
		Input:  variablesJson,
	}
	err = db.Create(row, dal.From(collector.table))
	if err != nil {
		collector.checkError(errors.Default.Wrap(err, `not created row table in graphql collector`))
		return
	}

	var results []interface{}
	if collector.args.ResponseParserWithDataErrors != nil {
		results, err = collector.args.ResponseParserWithDataErrors(query, variables, dataErrors)
	} else {
		results, err = collector.args.ResponseParser(query, variables)
	}
	if err != nil {
		if errors.Is(err, ErrFinishCollect) {
			logger.Info("collector finish by parser, rawId: #%d", row.ID)
			handler = nil
		} else {
			collector.checkError(errors.Default.Wrap(err, `not parsed response in graphql collector`))
			return
		}
	}
	err = collector.BatchSaveWithOrigin(divider, results, row)
	if err != nil {
		collector.checkError(err)
		return
	}

	collector.args.Ctx.IncProgress(1)
	if handler != nil {
		// trigger next fetch, but return if ErrFinishCollect got from ResponseParser
		err = handler(query)
		if err != nil {
			collector.checkError(errors.Default.Wrap(err, `handle failed in graphql collector`))
			return
		}
	}
}

func (collector *GraphqlCollector) checkError(err error) {
	if err == nil {
		return
	}
	collector.workerErrors = append(collector.workerErrors, err)
}

// HasError return if any error occurred
func (collector *GraphqlCollector) HasError() bool {
	return len(collector.workerErrors) > 0
}

var _ plugin.SubTask = (*GraphqlCollector)(nil)
