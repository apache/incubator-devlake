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
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// bridge to current implementation at this point
// TODO: implement another TaskContext for distributed runner/worker

// shared by TasContext and SubTaskContext
type defaultExecContext struct {
	cfg      *viper.Viper
	logger   core.Logger
	db       *gorm.DB
	ctx      context.Context
	name     string
	data     interface{}
	total    int
	current  int64
	mu       sync.Mutex
	progress chan core.RunningProgress
}

func newDefaultExecContext(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	data interface{},
	progress chan core.RunningProgress,
) *defaultExecContext {
	return &defaultExecContext{
		cfg:      cfg,
		logger:   logger,
		db:       db,
		ctx:      ctx,
		name:     name,
		data:     data,
		progress: progress,
	}
}

func (c *defaultExecContext) GetName() string {
	return c.name
}

func (c *defaultExecContext) GetConfig(name string) string {
	return c.cfg.GetString(name)
}

func (c *defaultExecContext) GetDb() *gorm.DB {
	return c.db
}

func (c *defaultExecContext) GetContext() context.Context {
	return c.ctx
}

func (c *defaultExecContext) GetData() interface{} {
	return c.data
}

func (c *defaultExecContext) GetLogger() core.Logger {
	return c.logger
}

func (c *defaultExecContext) SetProgress(progressType core.ProgressType, current int, total int) {
	c.current = int64(current)
	c.total = total

	if c.progress != nil {
		c.progress <- core.RunningProgress{
			Type:    progressType,
			Current: current,
			Total:   total,
		}
	}
}

func (c *defaultExecContext) IncProgress(progressType core.ProgressType, quantity int) {
	atomic.AddInt64(&c.current, int64(quantity))
	current := c.current
	if c.progress != nil {
		c.progress <- core.RunningProgress{
			Type:    progressType,
			Current: int(current),
			Total:   c.total,
		}
		// subtask progress may go too fast, remove old messages because they don't matter any more
		if progressType == core.SubTaskSetProgress {
			for len(c.progress) > 1 {
				<-c.progress
			}
		}
	}
}

func (c *defaultExecContext) fork(name string) *defaultExecContext {
	return newDefaultExecContext(
		c.cfg,
		c.logger.Nested(name),
		c.db,
		c.ctx,
		name,
		c.data,
		c.progress,
	)
}

// TaskContext default implementation
type DefaultTaskContext struct {
	*defaultExecContext
	subtasks    map[string]bool
	subtaskCtxs map[string]*DefaultSubTaskContext
}

func (c *DefaultTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(core.TaskSetProgress, current, total)
	c.logger.Info("total step: %d", c.total)
}

func (c *DefaultTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(core.TaskIncProgress, quantity)
	c.logger.Info("finished step: %d / %d", c.current, c.total)
}

// SubTaskContext default implementation
type DefaultSubTaskContext struct {
	*defaultExecContext
	taskCtx *DefaultTaskContext
}

func (c *DefaultSubTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(core.SubTaskSetProgress, current, total)
	if total > -1 {
		c.logger.Info("total records: %d", c.total)
	}
}

func (c *DefaultSubTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(core.SubTaskIncProgress, quantity)
	c.logger.Info("finished records: %d", c.current)
}

func NewDefaultTaskContext(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	subtasks map[string]bool,
	progress chan core.RunningProgress,
) core.TaskContext {
	return &DefaultTaskContext{
		newDefaultExecContext(cfg, logger, db, ctx, name, nil, progress),
		subtasks,
		make(map[string]*DefaultSubTaskContext),
	}
}

func (c *DefaultTaskContext) SubTaskContext(subtask string) (core.SubTaskContext, error) {
	// no need to lock at this point because subtasks is written only once
	if run, ok := c.subtasks[subtask]; ok {
		if run {
			// now, create a sub task context if it didn't exist
			c.mu.Lock()
			if c.subtaskCtxs[subtask] == nil {
				c.subtaskCtxs[subtask] = &DefaultSubTaskContext{
					c.defaultExecContext.fork(subtask),
					c,
				}
			}
			c.mu.Unlock()
			return c.subtaskCtxs[subtask], nil
		}
		// subtasks is skipped
		return nil, nil
	}
	// invalid subtask name
	return nil, fmt.Errorf("subtask %s doesn't exist", subtask)
}

// This returns a stand-alone core.SubTaskContext,
// not attached to any core.TaskContext.
// Use this if you need to run/debug a subtask without
// going through the usual workflow.
func NewStandaloneSubTaskContext(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	data interface{},
) core.SubTaskContext {
	return &DefaultSubTaskContext{
		newDefaultExecContext(cfg, logger, db, ctx, name, data, nil),
		nil,
	}
}

func (c *DefaultTaskContext) SetData(data interface{}) {
	c.data = data
}

var _ core.TaskContext = (*DefaultTaskContext)(nil)

func (c *DefaultSubTaskContext) TaskContext() core.TaskContext {
	return c.taskCtx
}

var _ core.SubTaskContext = (*DefaultSubTaskContext)(nil)
