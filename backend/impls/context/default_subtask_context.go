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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"time"
)

// DefaultSubTaskContext is default implementation
type DefaultSubTaskContext struct {
	*defaultExecContext
	taskCtx          *DefaultTaskContext
	LastProgressTime time.Time
}

// SetProgress FIXME ...
func (c *DefaultSubTaskContext) SetProgress(current int, total int) {
	c.defaultExecContext.SetProgress(plugin.SubTaskSetProgress, current, total)
	if total > -1 {
		c.BasicRes.GetLogger().Info("total jobs: %d", c.total)
	}
}

// IncProgress FIXME ...
func (c *DefaultSubTaskContext) IncProgress(quantity int) {
	c.defaultExecContext.IncProgress(plugin.SubTaskIncProgress, quantity)
	if c.LastProgressTime.IsZero() || c.LastProgressTime.Add(3*time.Second).Before(time.Now()) || c.current%1000 == 0 {
		c.LastProgressTime = time.Now()
		c.BasicRes.GetLogger().Info("finished records: %d", c.current)
	} else {
		c.BasicRes.GetLogger().Debug("finished records: %d", c.current)
	}
}

// TaskContext FIXME ...
func (c *DefaultSubTaskContext) TaskContext() plugin.TaskContext {
	if c.taskCtx == nil {
		return nil
	}
	return c.taskCtx
}

// NewStandaloneSubTaskContext returns a stand-alone plugin.SubTaskContext,
// not attached to any plugin.TaskContext.
// Use this if you need to run/debug a subtask without
// going through the usual workflow.
func NewStandaloneSubTaskContext(
	ctx gocontext.Context,
	basicRes context.BasicRes,
	name string,
	data interface{},
) plugin.SubTaskContext {
	return &DefaultSubTaskContext{
		newDefaultExecContext(ctx, basicRes, name, data, nil),
		nil,
		time.Time{},
	}
}

var _ plugin.SubTaskContext = (*DefaultSubTaskContext)(nil)
