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
	"time"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/plugins/core"
)

// TODO: move this file to `impl` module

// shared by TasContext and SubTaskContext
type defaultExecContext struct {
	core.BasicRes
	ctx      context.Context
	name     string
	data     interface{}
	total    int
	current  int64
	mu       sync.Mutex
	progress chan core.RunningProgress
}

func newDefaultExecContext(
	ctx context.Context,
	basicRes core.BasicRes,
	name string,
	data interface{},
	progress chan core.RunningProgress,
) *defaultExecContext {
	return &defaultExecContext{
		BasicRes: basicRes,
		ctx:      ctx,
		name:     name,
		data:     data,
		progress: progress,
	}
}

func (c *defaultExecContext) GetName() string {
	return c.name
}

func (c *defaultExecContext) GetContext() context.Context {
	return c.ctx
}

func (c *defaultExecContext) GetData() interface{} {
	return c.data
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
		c.ctx,
		c.BasicRes.NestedLogger(name),
		name,
		c.data,
		c.progress,
	)
}

// DefaultTaskContext is TaskContext default implementation
type DefaultTaskContext struct {
	*defaultExecContext
	subtasks    map[string]bool
	subtaskCtxs map[string]*DefaultSubTaskContext
}

// SetProgress FIXME ...
func (c *DefaultTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(core.TaskSetProgress, current, total)
	c.BasicRes.GetLogger().Info("total step: %d", c.total)
}

// IncProgress FIXME ...
func (c *DefaultTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(core.TaskIncProgress, quantity)
	c.BasicRes.GetLogger().Info("finished step: %d / %d", c.current, c.total)
}

// DefaultSubTaskContext is default implementation
type DefaultSubTaskContext struct {
	*defaultExecContext
	taskCtx          *DefaultTaskContext
	LastProgressTime time.Time
}

// SetProgress FIXME ...
func (c *DefaultSubTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(core.SubTaskSetProgress, current, total)
	if total > -1 {
		c.BasicRes.GetLogger().Info("total jobs: %d", c.total)
	}
}

// IncProgress FIXME ...
func (c *DefaultSubTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(core.SubTaskIncProgress, quantity)
	if c.LastProgressTime.IsZero() || c.LastProgressTime.Add(3*time.Second).Before(time.Now()) || c.current%1000 == 0 {
		c.LastProgressTime = time.Now()
		c.BasicRes.GetLogger().Info("finished records: %d", c.current)
	} else {
		c.BasicRes.GetLogger().Debug("finished records: %d", c.current)
	}
}

// NewDefaultTaskContext holds everything needed by the task execution.
func NewDefaultTaskContext(
	ctx context.Context,
	basicRes core.BasicRes,
	name string,
	subtasks map[string]bool,
	progress chan core.RunningProgress,
) core.TaskContext {
	return &DefaultTaskContext{
		newDefaultExecContext(ctx, basicRes, name, nil, progress),
		subtasks,
		make(map[string]*DefaultSubTaskContext),
	}
}

// SubTaskContext FIXME ...
func (c *DefaultTaskContext) SubTaskContext(subtask string) (core.SubTaskContext, errors.Error) {
	// no need to lock at this point because subtasks is written only once
	if run, ok := c.subtasks[subtask]; ok {
		if run {
			// now, create a sub task context if it didn't exist
			c.defaultExecContext.mu.Lock()
			if c.subtaskCtxs[subtask] == nil {
				c.subtaskCtxs[subtask] = &DefaultSubTaskContext{
					c.defaultExecContext.fork(subtask),
					c,
					time.Time{},
				}
			}
			c.defaultExecContext.mu.Unlock()
			return c.subtaskCtxs[subtask], nil
		}
		// subtasks is skipped
		return nil, nil
	}
	// invalid subtask name
	return nil, errors.Default.New(fmt.Sprintf("subtask %s doesn't exist", subtask))
}

// NewStandaloneSubTaskContext returns a stand-alone core.SubTaskContext,
// not attached to any core.TaskContext.
// Use this if you need to run/debug a subtask without
// going through the usual workflow.
func NewStandaloneSubTaskContext(
	ctx context.Context,
	basicRes core.BasicRes,
	name string,
	data interface{},
) core.SubTaskContext {
	return &DefaultSubTaskContext{
		newDefaultExecContext(ctx, basicRes, name, data, nil),
		nil,
		time.Time{},
	}
}

// SetData FIXME ...
func (c *DefaultTaskContext) SetData(data interface{}) {
	c.data = data
}

var _ core.TaskContext = (*DefaultTaskContext)(nil)

// TaskContext FIXME ...
func (c *DefaultSubTaskContext) TaskContext() core.TaskContext {
	if c.taskCtx == nil {
		return nil
	}
	return c.taskCtx
}

var _ core.SubTaskContext = (*DefaultSubTaskContext)(nil)
