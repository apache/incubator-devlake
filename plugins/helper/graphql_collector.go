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

package helper

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"net/http"
	"reflect"
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

type GraphqlQueryPageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type GraphqlAsyncResponseHandler func(res *http.Response) error

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
	// GetPageInfo is to tell `GraphqlCollector` is page information
	GetPageInfo    func(query interface{}, args *GraphqlCollectorArgs) (*GraphqlQueryPageInfo, error)
	BatchSize      int
	ResponseParser func(query interface{}, variables map[string]interface{}) ([]interface{}, error)
}

type GraphqlCollector struct {
	*RawDataSubTask
	args *GraphqlCollectorArgs
}

// NewGraphqlCollector allocates a new GraphqlCollector with the given args.
// GraphqlCollector can help us collect data from api with ease, pass in a AsyncGraphqlClient and tell it which part
// of response we want to save, GraphqlCollector will collect them from remote server and store them into database.
func NewGraphqlCollector(args GraphqlCollectorArgs) (*GraphqlCollector, error) {
	// process args
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to compile UrlTemplate: %w", err)
	}
	if args.GraphqlClient == nil {
		return nil, fmt.Errorf("ApiClient is required")
	}
	if args.ResponseParser == nil {
		return nil, fmt.Errorf("ResponseParser is required")
	}
	apicllector := &GraphqlCollector{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	if args.InputStep == 0 {
		args.InputStep = 1
	}
	//if args.AfterResponse != nil {
	//	apicllector.SetAfterResponse(args.AfterResponse)
	//} else {
	//	apicllector.SetAfterResponse(func(res *http.Response) error {
	//		if res.StatusCode == http.StatusUnauthorized {
	//			return fmt.Errorf("authentication failed, please check your AccessToken")
	//		}
	//		return nil
	//	})
	//}
	return apicllector, nil
}

// Start collection
func (collector *GraphqlCollector) Execute() error {
	logger := collector.args.Ctx.GetLogger()
	logger.Info("start graphql collection")

	// make sure table is created
	db := collector.args.Ctx.GetDal()
	err := db.AutoMigrate(&RawData{}, dal.From(collector.table))
	if err != nil {
		return err
	}

	// flush data if not incremental collection
	err = db.Delete(&RawData{}, dal.From(collector.table), dal.Where("params = ?", collector.params))
	if err != nil {
		return err
	}
	divider := NewBatchSaveDivider(collector.args.Ctx, collector.args.BatchSize, collector.table, collector.params)

	collector.args.Ctx.SetProgress(0, -1)
	if collector.args.Input != nil {
		iterator := collector.args.Input
		defer iterator.Close()
		apiClient := collector.args.GraphqlClient
		for iterator.HasNext() && !apiClient.HasError() {
			if collector.args.InputStep == 1 {
				input, err := iterator.Fetch()
				if err != nil {
					break
				}
				collector.exec(divider, input)
			} else {
				var inputs []interface{}
				for i := 0; i < collector.args.InputStep && iterator.HasNext(); i++ {
					input, err := iterator.Fetch()
					if err != nil {
						break
					}
					inputs = append(inputs, input)
				}
				collector.exec(divider, inputs)
			}
		}
	} else {
		// or we just did it once
		collector.exec(divider, nil)
	}

	logger.Debug("wait for all async api to finished")

	err = collector.args.GraphqlClient.Wait()
	if err != nil {
		logger.Info("end api collection error: %w", err)
	} else {
		logger.Info("end api collection without error")
	}
	err = divider.Close()

	return err
}

func (collector *GraphqlCollector) exec(divider *BatchSaveDivider, input interface{}) {
	inputJson, err := json.Marshal(input)
	if err != nil {
		panic(err)
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

// fetchPagesDetermined fetches data of all pages for APIs that return paging information
func (collector *GraphqlCollector) fetchOneByOne(divider *BatchSaveDivider, reqData *GraphqlRequestData) {
	// fetch first page
	var fetchNextPage func(query interface{}) error
	fetchNextPage = func(query interface{}) error {
		pageInfo, err := collector.args.GetPageInfo(query, collector.args)
		if err != nil {
			return fmt.Errorf("fetchPagesDetermined get totalPages faileds: %s", err.Error())
		}
		if pageInfo.HasNextPage {
			collector.args.GraphqlClient.NextTick(func() error {
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
			})
		}
		return nil
	}
	collector.fetchAsync(divider, reqData, fetchNextPage)
}

func (collector *GraphqlCollector) fetchAsync(divider *BatchSaveDivider, reqData *GraphqlRequestData, handler func(query interface{}) error) {
	if reqData.Pager == nil {
		reqData.Pager = &CursorPager{
			SkipCursor: nil,
			Size:       collector.args.PageSize,
		}
	}
	query, variables, err := collector.args.BuildQuery(reqData)
	if err != nil {
		panic(err)
	}

	logger := collector.args.Ctx.GetLogger()
	err = collector.args.GraphqlClient.Query(query, variables)
	if err != nil {
		panic(err)
	}
	defer logger.Debug("fetchAsync >>> done for %v %v", query, variables)

	paramsBytes, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}
	db := collector.args.Ctx.GetDal()
	row := &RawData{
		Params: collector.params,
		Data:   paramsBytes,
		Url:    `TODO query`,
		Input:  reqData.InputJSON,
	}
	err = db.Create(row, dal.From(collector.table))
	if err != nil {
		panic(err)
	}

	results, err := collector.args.ResponseParser(query, variables)
	if err != nil {
		panic(err)
	}

	RAW_DATA_ORIGIN := "RawDataOrigin"
	// batch save divider
	for _, result := range results {
		// get the batch operator for the specific type
		batch, err := divider.ForType(reflect.TypeOf(result))
		if err != nil {
			panic(err)
		}
		// set raw data origin field
		origin := reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN)
		if origin.IsValid() {
			origin.Set(reflect.ValueOf(common.RawDataOrigin{
				RawDataTable:  collector.table,
				RawDataId:     row.ID,
				RawDataParams: row.Params,
			}))
		}
		// records get saved into db when slots were max outed
		err = batch.Add(result)
		if err != nil {
			panic(err)
		}
		collector.args.Ctx.IncProgress(1)
	}
	if err != nil {
		panic(err)
	}
	collector.args.Ctx.IncProgress(1)
	if handler != nil {
		err = handler(query)
		if err != nil {
			panic(err)
		}
	}
}

var _ core.SubTask = (*ApiCollector)(nil)
