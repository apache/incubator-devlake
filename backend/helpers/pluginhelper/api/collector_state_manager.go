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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

// CollectorStateManager manages the state of the collector. It is used to determine whether
// the collector should run in incremental mode or full sync mode and what time range to collect.
type CollectorStateManager struct {
	IsIncremental bool
	Since         *time.Time
	Until         *time.Time
}

// type CollectorOptions struct {
// 	TimeAfter string `json:"timeAfter,omitempty" mapstructure:"timeAfter,omitempty"`
// }

// NewCollectorStateManager create a new CollectorStateManager
func NewCollectorStateManager(subtaskCtx plugin.SubTaskContext, rawTable, rawParams string) (*CollectorStateManager, errors.Error) {
	db := subtaskCtx.GetDal()
	syncPolicy := subtaskCtx.TaskContext().SyncPolicy()

	// CollectorLatestState retrieves the latest collector state from the database
	oldState := models.CollectorLatestState{}
	err := db.First(&oldState, dal.Where(`raw_data_table = ? AND raw_data_params = ?`, rawTable, rawParams))
	if err != nil {
		if db.IsErrorNotFound(err) {
			oldState = models.CollectorLatestState{
				RawDataTable:  rawTable,
				RawDataParams: rawParams,
			}
		} else {
			return nil, errors.Default.Wrap(err, "failed to load JiraLatestCollectorMeta")
		}
	}
	// Extract timeAfter and latestSuccessStart from old state
	oldTimeAfter := oldState.TimeAfter
	oldLatestSuccessStart := oldState.LatestSuccessStart

	// Calculate incremental and since based on syncPolicy and old state
	var isIncremental bool
	var since *time.Time

	if oldLatestSuccessStart == nil {
		// 1. If no oldState.LatestSuccessStart, not incremental and since is syncPolicy.TimeAfter
		isIncremental = false
		if syncPolicy != nil {
			since = syncPolicy.TimeAfter
		}
	} else if syncPolicy == nil {
		// 2. If no syncPolicy, incremental and since is oldState.LatestSuccessStart
		isIncremental = true
		since = oldLatestSuccessStart
	} else if syncPolicy.FullSync {
		// 3. If fullSync true, not incremental and since is syncPolicy.TimeAfter
		isIncremental = false
		since = syncPolicy.TimeAfter
	} else if syncPolicy.TimeAfter == nil {
		// 4. If no syncPolicy TimeAfter, incremental and since is oldState.LatestSuccessStart
		isIncremental = true
		since = oldLatestSuccessStart
	} else {
		// 5. If syncPolicy.TimeAfter not nil
		if oldTimeAfter != nil && syncPolicy.TimeAfter.Before(*oldTimeAfter) {
			// 4.1 If oldTimeAfter not nil and syncPolicy.TimeAfter before oldTimeAfter, incremental is false and since is syncPolicy.TimeAfter
			isIncremental = false
			since = syncPolicy.TimeAfter
		} else {
			// 4.2 If oldTimeAfter nil or syncPolicy.TimeAfter after oldTimeAfter, incremental is true and since is oldState.LatestSuccessStart
			isIncremental = true
			since = oldLatestSuccessStart
		}
	}

	currentTime := time.Now()
	oldState.LatestSuccessStart = &currentTime
	oldState.TimeAfter = syncPolicy.TimeAfter

	return &CollectorStateManager{
		IsIncremental: isIncremental,
		Since:         since,
		Until:         &currentTime,
	}, nil
}

func (c *CollectorStateManager) Save() errors.Error {
	return nil
}
