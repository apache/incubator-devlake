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
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// StatefulApiCollector runs multiple collectors as a single subtask and maintains the state of the collector
// mainly the time range to collect across multiple collections. It is useful when you need to support timeAfter
// and diff sync for APIs that do not support filtering by the updated date.
type StatefulApiCollector struct {
	RawDataSubTaskArgs
	CollectorStateManager
	// *ApiCollector
	// *GraphqlCollector
	nestedCollectors []plugin.SubTask
}

// NewStatefulApiCollector create a new StatefulApiCollector
func NewStatefulApiCollector(args RawDataSubTaskArgs) (*StatefulApiCollector, errors.Error) {
	syncPolicy := args.Ctx.TaskContext().SyncPolicy()
	rawDataSubTask, err := NewRawDataSubTask(args)
	if err != nil {
		return nil, err
	}
	stateManager, err := NewCollectorStateManager(args.Ctx, syncPolicy, rawDataSubTask.table, rawDataSubTask.params)
	if err != nil {
		return nil, err
	}
	return &StatefulApiCollector{
		RawDataSubTaskArgs:    args,
		CollectorStateManager: *stateManager,
	}, nil
}

// InitCollector appends a new collector to the list
func (m *StatefulApiCollector) InitCollector(args ApiCollectorArgs) errors.Error {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	args.Incremental = m.CollectorStateManager.IsIncremental()
	apiCollector, err := NewApiCollector(args)
	if err != nil {
		return err
	}
	m.nestedCollectors = append(m.nestedCollectors, apiCollector)
	return nil
}

// InitGraphQLCollector appends a new GraphQL collector to the list
func (m *StatefulApiCollector) InitGraphQLCollector(args GraphqlCollectorArgs) errors.Error {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	args.Incremental = m.CollectorStateManager.IsIncremental()
	graphqlCollector, err := NewGraphqlCollector(args)
	if err != nil {
		return err
	}
	m.nestedCollectors = append(m.nestedCollectors, graphqlCollector)
	return nil
}

// Execute all nested collectors and save the state if all collectors succeed
func (m *StatefulApiCollector) Execute() errors.Error {
	for _, subtask := range m.nestedCollectors {
		err := subtask.Execute()
		if err != nil {
			return err
		}
	}

	return m.CollectorStateManager.Close()
}

// NewStatefulApiCollectorForFinalizableEntity aims to add timeFilter/diffSync support for
// APIs that do NOT support filtering data by the updated date. However, it comes with the
// following constraints:
//  1. The entity is a short-lived object or it is likely to be irrelevant
//     a. ci/id pipelines are short-lived objects
//     b. pull request might took a year to be closed or never, but it is likely irrelevant
//  2. The entity must be Finalizable, meaning no future modifications will happen to it once it
//     enter some sort of `Closed`/`Finished` status.
//  3. The API must fit one of the following traits:
//     a. it supports filtering by Created Date, in this case, you must implement the filtering
//     via the `UrlTemplate`, `Query` or `Header` hook based on the API specification.
//     b. or sorting by Created Date in Descending order, in this case, you must use `Concurrency`
//     or `GetNextPageCustomData` instead of `GetTotalPages` for Undetermined Strategy since we have
//     to stop the process in the middle.
//
// Assuming the API fits the bill, the strategies can be categoried into:
//   - Determined Strategy: if the API supports filtering by the Created Date, use the `GetTotalPages` hook
//   - Undetermind Strategy: if the API supports sorting by the Created Date in Descending order and
//     fetching by Page Number, use the `Concurrent` hook
//   - Sequential Strategy: if the API supports sorting by the Created Date in Descending order but
//     the next page can only be fetched by the Cursor/Token from the previous page, use the `GetNextPageCustomData` hook
func NewStatefulApiCollectorForFinalizableEntity(args FinalizableApiCollectorArgs) (plugin.SubTask, errors.Error) {
	// create a manager which could execute multiple collector but acts as a single subtask to callers
	manager, err := NewStatefulApiCollector(RawDataSubTaskArgs{
		Ctx:     args.Ctx,
		Options: args.Options,
		Params:  args.Params,
		Table:   args.Table,
	})
	if err != nil {
		return nil, err
	}

	createdAfter := manager.CollectorStateManager.GetSince()
	isIncremental := manager.CollectorStateManager.IsIncremental()

	// step 1: create a collector to collect newly added records
	err = manager.InitCollector(ApiCollectorArgs{
		ApiClient: args.ApiClient,
		// common
		Incremental: isIncremental,
		UrlTemplate: args.CollectNewRecordsByList.UrlTemplate,
		Query: func(reqData *RequestData) (url.Values, errors.Error) {
			if args.CollectNewRecordsByList.Query != nil {
				return args.CollectNewRecordsByList.Query(reqData, createdAfter)
			}
			return nil, nil
		},
		Header: func(reqData *RequestData) (http.Header, errors.Error) {
			if args.CollectNewRecordsByList.Header != nil {
				return args.CollectNewRecordsByList.Header(reqData, createdAfter)
			}
			return nil, nil
		},
		MinTickInterval: args.CollectNewRecordsByList.MinTickInterval,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			items, err := args.CollectNewRecordsByList.ResponseParser(res)
			if err != nil {
				return nil, err
			}
			if len(items) == 0 {
				return nil, nil
			}

			// time filter or diff sync
			if createdAfter != nil && args.CollectNewRecordsByList.GetCreated != nil {
				// if the first record of the page was created before createdAfter, return emtpy set and stop
				firstCreated, err := args.CollectNewRecordsByList.GetCreated(items[0])
				if err != nil {
					return nil, err
				}
				if firstCreated.Before(*createdAfter) {
					return nil, ErrFinishCollect
				}
				// if the last record was created before createdAfter, return records and stop
				lastCreated, err := args.CollectNewRecordsByList.GetCreated(items[len(items)-1])
				if err != nil {
					return nil, err
				}
				if lastCreated.Before(*createdAfter) {
					return items, ErrFinishCollect
				}
			}
			return items, err
		},
		AfterResponse: args.CollectNewRecordsByList.AfterResponse,
		RequestBody:   args.CollectNewRecordsByList.RequestBody,
		Method:        args.CollectNewRecordsByList.Method,
		// pagination
		PageSize:              args.CollectNewRecordsByList.PageSize,
		Concurrency:           args.CollectNewRecordsByList.Concurrency,
		GetNextPageCustomData: args.CollectNewRecordsByList.GetNextPageCustomData,
		GetTotalPages:         args.CollectNewRecordsByList.GetTotalPages,
	})

	if err != nil {
		return nil, err
	}

	syncPolicy := args.Ctx.TaskContext().SyncPolicy()
	if args.CollectUnfinishedDetails == nil || (syncPolicy != nil && syncPolicy.FullSync) {
		return manager, nil
	}

	// step 2: create another collector to collect updated records
	// TODO: this creates cursor before previous step gets executed, which is too early, to be optimized
	input, err := args.CollectUnfinishedDetails.BuildInputIterator()
	if err != nil {
		return nil, err
	}
	err = manager.InitCollector(ApiCollectorArgs{
		ApiClient: args.ApiClient,
		// common
		Incremental: true,
		Input:       input,
		UrlTemplate: args.CollectUnfinishedDetails.UrlTemplate,
		Query: func(reqData *RequestData) (url.Values, errors.Error) {
			if args.CollectUnfinishedDetails.Query != nil {
				return args.CollectUnfinishedDetails.Query(reqData, createdAfter)
			}
			return nil, nil
		},
		Header: func(reqData *RequestData) (http.Header, errors.Error) {
			if args.CollectUnfinishedDetails.Header != nil {
				return args.CollectUnfinishedDetails.Header(reqData, createdAfter)
			}
			return nil, nil
		},
		MinTickInterval: args.CollectUnfinishedDetails.MinTickInterval,
		ResponseParser:  args.CollectUnfinishedDetails.ResponseParser,
		AfterResponse:   args.CollectUnfinishedDetails.AfterResponse,
		RequestBody:     args.CollectUnfinishedDetails.RequestBody,
		Method:          args.CollectUnfinishedDetails.Method,
	})
	return manager, err
}

type FinalizableApiCollectorArgs struct {
	RawDataSubTaskArgs
	ApiClient                RateLimitedApiClient
	CollectNewRecordsByList  FinalizableApiCollectorListArgs
	CollectUnfinishedDetails *FinalizableApiCollectorDetailArgs
}

// FinalizableApiCollectorCommonArgs is the common arguments for both list and detail collectors
// Note that all request-related arguments would be called or utilized before any response-related arguments
type FinalizableApiCollectorCommonArgs struct {
	UrlTemplate     string                                                                          // required, url path template for the request, e.g. repos/{{ .Params.Name }}/pulls or incident/{{ .Input.Number }} (if using iterators)
	Method          string                                                                          // optional, request method, e.g. GET(default), POST, PUT, DELETE
	Query           func(reqData *RequestData, createdAfter *time.Time) (url.Values, errors.Error)  // optional, build query params for the request
	Header          func(reqData *RequestData, createdAfter *time.Time) (http.Header, errors.Error) // optional, build header for the request
	RequestBody     func(reqData *RequestData) map[string]interface{}                               // optional, build request body for the request if the Method set to POST or PUT
	MinTickInterval *time.Duration                                                                  // optional, minimum interval between two requests, some endpoints might have a more conservative rate limit than others within the same instance, you can mitigate this by setting a higher MinTickInterval to override the connection level rate limit.
	AfterResponse   plugin.ApiClientAfterResponse                                                   // optional, hook to run after each response, would be called before the ResponseParser
	ResponseParser  func(res *http.Response) ([]json.RawMessage, errors.Error)                      // required, parse the response body and return a list of entities
}

// FinalizableApiCollectorListArgs is the arguments for the list collector
type FinalizableApiCollectorListArgs struct {
	FinalizableApiCollectorCommonArgs
	GetCreated            func(item json.RawMessage) (time.Time, errors.Error)                                        // optional, to extract create date from a raw json of a single record, leave it be `nil` if API supports filtering by updated date (Don't forget to set the Query)
	PageSize              int                                                                                         // required, number of records per page
	Concurrency           int                                                                                         // required for Undetermined Strategy, number of concurrent requests
	GetNextPageCustomData func(prevReqData *RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) // required for Sequential Strategy, to extract the next page cursor from the given response
	GetTotalPages         func(res *http.Response, args *ApiCollectorArgs) (int, errors.Error)                        // required for Determined Strategy, to extract the total number of pages from the given response
}

// FinalizableApiCollectorDetailArgs is the arguments for the detail collector
type FinalizableApiCollectorDetailArgs struct {
	FinalizableApiCollectorCommonArgs
	BuildInputIterator func() (Iterator, errors.Error) // required, create an iterator that iterates through all unfinalized records in the database. These records will be fed as the "Input" (or {{ .Input.* }} in URLTemplate) argument back into FinalizableApiCollectorCommonArgs which makes the API calls to re-collect their newest states.
}
