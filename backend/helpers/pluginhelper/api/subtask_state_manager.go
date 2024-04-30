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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
)

// SubtaskCommonArgs is a struct that contains the common arguments for a subtask
type SubtaskCommonArgs struct {
	plugin.SubTaskContext
	Table         string // raw table name
	Params        any    // for filtering rows belonging to the scope (jira board, github repo) of the subtask
	SubtaskConfig any    // for determining whether the subtask should run in incremental or full sync mode
	BatchSize     int    // batch size for saving data
}

func (args *SubtaskCommonArgs) GetRawDataTable() string {
	return fmt.Sprintf("_raw_%s", args.Table)
}

func (args *SubtaskCommonArgs) GetRawDataParams() string {
	if args.Params == nil || reflect.ValueOf(args.Params).IsZero() {
		panic(errors.Default.New("Params is nil"))
	}
	return utils.ToJsonString(args.Params)
}

func (args *SubtaskCommonArgs) GetSubtaskConfig() string {
	return utils.ToJsonString(args.SubtaskConfig)
}

func (args *SubtaskCommonArgs) GetBatchSize() int {
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return args.BatchSize
}

// SubtaskStateManager manages the state of a subtask. It is used to determine whether
// the subtask should run in incremental mode or full sync mode and what time range to collect.
type SubtaskStateManager struct {
	db            dal.Dal
	state         *models.SubtaskState
	syncPolicy    *models.SyncPolicy
	isIncremental bool       // tells if the subtask should run in incremental mode or full sync mode
	since         *time.Time // the start time of the time range to work on
	until         *time.Time // the end time of the time range to work on
	config        string     // current configuration of the subtask for determining if the subtask should run in incremental or full sync mode
}

// NewSubtaskStateManager create a new SubtaskStateManager
func NewSubtaskStateManager(args *SubtaskCommonArgs) (stateManager *SubtaskStateManager, err errors.Error) {
	db := args.GetDal()
	syncPolicy := args.SubTaskContext.TaskContext().SyncPolicy()
	plugin := args.SubTaskContext.TaskContext().GetName()
	subtask := args.SubTaskContext.GetName()
	// load sync policy and make sure it is not nil
	if syncPolicy == nil {
		syncPolicy = &models.SyncPolicy{}
	}
	params := args.GetRawDataParams()
	// load the previous state from the database
	state := &models.SubtaskState{}
	err = db.First(state, dal.Where(`plugin = ? AND subtask =? AND params = ?`, plugin, subtask, params))
	if err != nil {
		if db.IsErrorNotFound(err) {
			state = &models.SubtaskState{
				Plugin:  plugin,
				Subtask: subtask,
				Params:  params,
			}
			err = nil
		} else {
			err = errors.Default.Wrap(err, "failed to load the previous subtask state")
			return
		}
	}
	// fullsync by default
	now := time.Now()
	stateManager = &SubtaskStateManager{
		db:            db,
		state:         state,
		syncPolicy:    syncPolicy,
		isIncremental: false,
		since:         syncPolicy.TimeAfter,
		until:         &now,
		config:        utils.ToJsonString(args.SubtaskConfig),
	}
	// fallback to the previous timeAfter if no new value
	if stateManager.since == nil {
		stateManager.since = state.TimeAfter
	}
	// if fullsync is set or no previous success start time, we are in the full sync mode
	if syncPolicy.FullSync || state.PrevStartedAt == nil {
		return
	}
	// if timeAfter is not set or NOT before the previous vaule, we are in the incremental mode
	if (syncPolicy.TimeAfter == nil || state.TimeAfter == nil || !syncPolicy.TimeAfter.Before(*state.TimeAfter)) &&
		// and the previous config is the same as the current config
		(state.PrevConfig == "" || state.PrevConfig == stateManager.config) {
		stateManager.isIncremental = true
		stateManager.since = state.PrevStartedAt
	}
	return
}

func (c *SubtaskStateManager) IsIncremental() bool {
	return c.isIncremental
}

func (c *SubtaskStateManager) GetSince() *time.Time {
	return c.since
}

func (c *SubtaskStateManager) GetUntil() *time.Time {
	return c.until
}

func (c *SubtaskStateManager) Close() errors.Error {
	// update timeAfter in the database only for fullsync mode
	if !c.isIncremental {
		// prefer non-nil value
		if c.syncPolicy.TimeAfter != nil {
			c.state.TimeAfter = c.syncPolicy.TimeAfter
		}
	}
	// always update the latest success start time
	c.state.PrevStartedAt = c.until
	c.state.PrevConfig = c.config
	return c.db.CreateOrUpdate(c.state)
}
