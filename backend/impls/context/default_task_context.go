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

package context

import (
	gocontext "context"
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"time"
)

// DefaultTaskContext is TaskContext default implementation
type DefaultTaskContext struct {
	*defaultExecContext
	subtasks    map[string]bool
	subtaskCtxs map[string]*DefaultSubTaskContext
}

// SetProgress FIXME ...
func (c *DefaultTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(plugin.TaskSetProgress, current, total)
	c.BasicRes.GetLogger().Info("total step: %d", c.total)
}

// IncProgress FIXME ...
func (c *DefaultTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(plugin.TaskIncProgress, quantity)
	c.BasicRes.GetLogger().Info("finished step: %d / %d", c.current, c.total)
}

// SubTaskContext FIXME ...
func (c *DefaultTaskContext) SubTaskContext(subtask string) (plugin.SubTaskContext, errors.Error) {
	// no need to lock at this point because subtasks is written only once
	if run, ok := c.subtasks[subtask]; ok {
		if run {
			// now, create a subtask context if it didn't exist
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

// SetData FIXME ...
func (c *DefaultTaskContext) SetData(data interface{}) {
	c.data = data
}

// NewDefaultTaskContext holds everything needed by the task execution.
func NewDefaultTaskContext(
	ctx gocontext.Context,
	basicRes context.BasicRes,
	name string,
	subtasks map[string]bool,
	progress chan plugin.RunningProgress,
) plugin.TaskContext {
	return &DefaultTaskContext{
		newDefaultExecContext(ctx, basicRes, name, nil, progress),
		subtasks,
		make(map[string]*DefaultSubTaskContext),
	}
}

var _ plugin.TaskContext = (*DefaultTaskContext)(nil)
