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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// ApiCollectorStateManager save collector state in framework table
type ApiCollectorStateManager struct {
	RawDataSubTaskArgs
	*ApiCollector
	*GraphqlCollector
	LatestState models.CollectorLatestState
	// Deprecating
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
		CreatedDateAfter:   createdDateAfter,
		TimeAfter:          timeAfter,
		ExecuteStart:       time.Now(),
	}, nil
}

// NewApiCollectorWithState create a new ApiCollectorStateManager
func NewApiCollectorWithState(args RawDataSubTaskArgs, createdDateAfter *time.Time) (*ApiCollectorStateManager, errors.Error) {
	return NewApiCollectorWithStateEx(args, createdDateAfter, nil)
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
	// fallback to CreatedDateAfter: collector should filter data by `created_date`
	return m.LatestState.CreatedDateAfter == nil || m.CreatedDateAfter != nil && !m.CreatedDateAfter.Before(*m.LatestState.CreatedDateAfter)
}

// InitCollector init the embedded collector
func (m *ApiCollectorStateManager) InitCollector(args ApiCollectorArgs) (err errors.Error) {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	m.ApiCollector, err = NewApiCollector(args)
	return err
}

// InitGraphQLCollector init the embedded collector
func (m *ApiCollectorStateManager) InitGraphQLCollector(args GraphqlCollectorArgs) (err errors.Error) {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	m.GraphqlCollector, err = NewGraphqlCollector(args)
	return err
}

// Execute the embedded collector and record execute state
func (m ApiCollectorStateManager) Execute() errors.Error {
	err := m.ApiCollector.Execute()
	if err != nil {
		return err
	}

	return m.updateState()
}

// ExecuteGraphQL the embedded collector and record execute state
func (m ApiCollectorStateManager) ExecuteGraphQL() errors.Error {
	err := m.GraphqlCollector.Execute()
	if err != nil {
		return err
	}

	return m.updateState()
}

func (m ApiCollectorStateManager) updateState() errors.Error {
	db := m.Ctx.GetDal()
	m.LatestState.LatestSuccessStart = &m.ExecuteStart
	m.LatestState.CreatedDateAfter = m.CreatedDateAfter
	m.LatestState.TimeAfter = m.TimeAfter
	return db.CreateOrUpdate(&m.LatestState)
}
