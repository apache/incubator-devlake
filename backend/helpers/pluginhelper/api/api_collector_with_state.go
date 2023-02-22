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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/common"
)

// ApiCollectorStateManager save collector state in framework table
type ApiCollectorStateManager struct {
	RawDataSubTaskArgs
	// *ApiCollector
	// *GraphqlCollector
	subtasks    []plugin.SubTask
	LatestState models.CollectorLatestState
	// Deprecating(timeAfter): to be deleted
	CreatedDateAfter *time.Time
	TimeAfter        *time.Time
	ExecuteStart     time.Time
}

// NewApiCollectorWithStateEx create a new ApiCollectorStateManager
func NewApiCollectorWithStateEx(args RawDataSubTaskArgs, createdDateAfter *time.Time, timeAfter *time.Time) (*ApiCollectorStateManager, errors.Error) {
	db := args.Ctx.GetDal()

	rawDataSubTask, err := NewRawDataSubTask(args)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Couldn't resolve raw subtask args")
	}
	latestState := models.CollectorLatestState{}
	err = db.First(&latestState, dal.Where(`raw_data_table = ? AND raw_data_params = ?`, rawDataSubTask.table, rawDataSubTask.params))
	if err != nil {
		if db.IsErrorNotFound(err) {
			latestState = models.CollectorLatestState{
				RawDataTable:  rawDataSubTask.table,
				RawDataParams: rawDataSubTask.params,
			}
		} else {
			return nil, errors.Default.Wrap(err, "failed to load JiraLatestCollectorMeta")
		}
	}
	return &ApiCollectorStateManager{
		RawDataSubTaskArgs: args,
		LatestState:        latestState,
		// Deprecating(timeAfter): to be deleted
		CreatedDateAfter: createdDateAfter,
		TimeAfter:        timeAfter,
		ExecuteStart:     time.Now(),
	}, nil
}

// NewApiCollectorWithState create a new ApiCollectorStateManager
// Deprecating(timeAfter): use NewStatefulApiCollector instead
func NewApiCollectorWithState(args RawDataSubTaskArgs, createdDateAfter *time.Time) (*ApiCollectorStateManager, errors.Error) {
	return NewApiCollectorWithStateEx(args, createdDateAfter, nil)
}

// NewApiCollectorWithState create a new ApiCollectorStateManager
func NewStatefulApiCollector(args RawDataSubTaskArgs, timeAfter *time.Time) (*ApiCollectorStateManager, errors.Error) {
	return NewApiCollectorWithStateEx(args, nil, timeAfter)
}

// IsIncremental indicates if the collector should operate in incremental mode
func (m *ApiCollectorStateManager) IsIncremental() bool {
	// the initial collection
	if m.LatestState.LatestSuccessStart == nil {
		return false
	}
	// prioritize TimeAfter parameter: collector should filter data by `updated_date`
	if m.TimeAfter != nil {
		return m.LatestState.TimeAfter == nil || !m.TimeAfter.Before(*m.LatestState.TimeAfter)
	}
	// Deprecating(timeAfter): to be removed
	// fallback to CreatedDateAfter: collector should filter data by `created_date`
	return m.LatestState.CreatedDateAfter == nil || m.CreatedDateAfter != nil && !m.CreatedDateAfter.Before(*m.LatestState.CreatedDateAfter)
}

// InitCollector init the embedded collector
func (m *ApiCollectorStateManager) InitCollector(args ApiCollectorArgs) errors.Error {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	apiCollector, err := NewApiCollector(args)
	if err != nil {
		return err
	}
	m.subtasks = append(m.subtasks, apiCollector)
	return nil
}

// InitGraphQLCollector init the embedded collector
func (m *ApiCollectorStateManager) InitGraphQLCollector(args GraphqlCollectorArgs) errors.Error {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	graphqlCollector, err := NewGraphqlCollector(args)
	if err != nil {
		return err
	}
	m.subtasks = append(m.subtasks, graphqlCollector)
	return nil
}

// Execute the embedded collector and record execute state
func (m *ApiCollectorStateManager) Execute() errors.Error {
	for _, subtask := range m.subtasks {
		err := subtask.Execute()
		if err != nil {
			return err
		}
	}

	db := m.Ctx.GetDal()
	m.LatestState.LatestSuccessStart = &m.ExecuteStart
	// Deprecating(timeAfter): to be deleted
	m.LatestState.CreatedDateAfter = m.CreatedDateAfter
	m.LatestState.TimeAfter = m.TimeAfter
	return db.CreateOrUpdate(&m.LatestState)
}

// NewStatefulApiCollectorForFinalizableEntity aims to add timeFilter/diffSync support for
// APIs that do NOT support filtering data by updated date. However, it comes with the
// following constraints:
//  1. The entity is a short-lived object or it is likely to be irrelevant
//     a. ci/id pipelines are short-lived objects
//     b. pull request might took a year to be closed or never, but it is likely irrelevant
//  2. The entity must be Finalizable: when it is finalized, no modification forever
//  3. The API must fit one of the following traits:
//     a. it supports filtering by Created Date, in this case, you may specify the `GetTotalPages`
//     option to fetch data with Determined Strategy if possible.
//     b. or sorting by Created Date in Descending order, this this case, you must use `Concurrency`
//     or `GetNextPageCustomData` instead of `GetTotalPages` for Undetermined Strategy since we have
//     to stop the process in the middle.
func NewStatefulApiCollectorForFinalizableEntity(args FinalizableApiCollectorArgs) (plugin.SubTask, errors.Error) {
	// create a manager which could execute multiple collector but acts as a single subtask to callers
	manager, err := NewStatefulApiCollector(RawDataSubTaskArgs{
		Ctx:    args.Ctx,
		Params: args.Params,
		Table:  args.Table,
	}, args.TimeAfter)
	if err != nil {
		return nil, err
	}

	// // prepare the basic variables
	var isIncremental = manager.IsIncremental()
	var createdAfter *time.Time
	if isIncremental {
		createdAfter = manager.LatestState.LatestSuccessStart
	} else {
		createdAfter = manager.TimeAfter
	}

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
		// GetTotalPages:         args.CollectNewRecordsByList.GetTotalPages,
	})

	if err != nil {
		return nil, err
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
	TimeAfter                *time.Time // leave it be nil to disable time filter
	CollectNewRecordsByList  FinalizableApiCollectorListArgs
	CollectUnfinishedDetails FinalizableApiCollectorDetailArgs
}

type FinalizableApiCollectorCommonArgs struct {
	UrlTemplate     string `comment:"GoTemplate for API url"`
	Query           func(reqData *RequestData, createdAfter *time.Time) (url.Values, errors.Error)
	Header          func(reqData *RequestData, createdAfter *time.Time) (http.Header, errors.Error)
	MinTickInterval *time.Duration
	ResponseParser  func(res *http.Response) ([]json.RawMessage, errors.Error)
	AfterResponse   common.ApiClientAfterResponse
	RequestBody     func(reqData *RequestData) map[string]interface{}
	Method          string
}
type FinalizableApiCollectorListArgs struct {
	// optional, leave it be `nil` if API supports filtering by created date (Don't forget to set the Query)
	GetCreated func(item json.RawMessage) (time.Time, errors.Error)
	FinalizableApiCollectorCommonArgs
	Concurrency           int
	PageSize              int
	GetNextPageCustomData func(prevReqData *RequestData, prevPageResponse *http.Response) (interface{}, errors.Error)
	// need to consider the data missing problem: what if new data gets created during collection?
	// GetTotalPages         func(res *http.Response, args *ApiCollectorArgs) (int, errors.Error)
}
type FinalizableApiCollectorDetailArgs struct {
	FinalizableApiCollectorCommonArgs
	BuildInputIterator func() (Iterator, errors.Error)
}
