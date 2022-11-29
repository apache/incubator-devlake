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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"gorm.io/gorm"
	"time"
)

type ApiCollectorStateManager struct {
	RawDataSubTaskArgs
	*ApiCollector
	db           dal.Dal
	LatestState  models.CollectorLatestState
	StartFrom    *time.Time
	ExecuteStart time.Time
}

func NewApiCollectorWithState(args RawDataSubTaskArgs, startFrom *time.Time) (*ApiCollectorStateManager, errors.Error) {
	db := args.Ctx.GetDal()

	rawDataSubTask, err := NewRawDataSubTask(args)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Couldn't resolve raw subtask args")
	}
	latestState := models.CollectorLatestState{}
	err = db.First(&latestState, dal.Where(`raw_data_table = ? AND raw_data_params = ?`, rawDataSubTask.table, rawDataSubTask.params))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		StartFrom:          startFrom,
		ExecuteStart:       time.Now(),
	}, nil
}

func (m ApiCollectorStateManager) CanIncrementCollect() bool {
	return m.LatestState.StartFrom != nil && m.LatestState.StartFrom.Equal(*m.StartFrom)
}

func (m *ApiCollectorStateManager) InitCollector(args ApiCollectorArgs) (err errors.Error) {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	m.ApiCollector, err = NewApiCollector(args)
	return err
}

func (m ApiCollectorStateManager) Execute() errors.Error {
	executeErr := m.ApiCollector.Execute()

	db := m.Ctx.GetDal()
	m.LatestState.LatestSuccessStart = &m.ExecuteStart
	m.LatestState.StartFrom = m.StartFrom
	saveErr := db.CreateOrUpdate(&m.LatestState)
	if saveErr != nil {
		if executeErr != nil {
			return errors.Default.Combine([]error{executeErr, saveErr})

		} else {
			return errors.Default.Wrap(saveErr, "error on saving JiraLatestCollectorMeta")
		}
	}
	return executeErr
}
