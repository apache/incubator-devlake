package helper

import (
	"context"
	"fmt"
	"sync"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
)

// bridge to current implementation at this point
// TODO: implement another TaskContext for distributed runner/worker

// shared by TasContext and SubTaskContext
type defaultExecContext struct {
	name    string
	ctx     context.Context
	data    interface{}
	logger  core.Logger
	total   int
	current int
	mu      sync.Mutex
}

func newDefaultExecContext(
	name string,
	ctx context.Context,
	data interface{},
	logger core.Logger,
) *defaultExecContext {
	return &defaultExecContext{
		name:   name,
		ctx:    ctx,
		data:   data,
		logger: logger,
	}
}

func (c *defaultExecContext) GetName() string {
	return c.name
}

func (c *defaultExecContext) GetConfig(name string) string {
	return config.GetConfig().GetString(name)
}

func (c *defaultExecContext) GetDb() *gorm.DB {
	return models.Db
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

func (c *defaultExecContext) SetProgress(current int, total int) {
	c.mu.Lock()
	c.current = current
	c.total = total
	c.mu.Unlock()
	if c.total >= 0 {
		c.logger.Info("set task %s progress: %d/%d", c.name, c.current, c.total)
	} else {
		c.logger.Info("set task %s progress: %d", c.name, c.current)
	}
}

func (c *defaultExecContext) IncProgress(quantity int) {
	c.mu.Lock()
	c.current += quantity
	c.mu.Unlock()
	if c.total >= 0 {
		c.logger.Info("increased task %s progress %d/%d", c.name, c.current, c.total)
	} else {
		c.logger.Info("increased task %s progress %d", c.name, c.current)
	}
}

func (c *defaultExecContext) fork(name string) *defaultExecContext {
	return newDefaultExecContext(name, c.ctx, c.data, c.logger.Nested(name))
}

// TaskContext default implementation
type DefaultTaskContext struct {
	*defaultExecContext
	subtasks    map[string]bool
	subtaskCtxs map[string]*DefaultSubTaskContext
}

// SubTaskContext default implementation
type DefaultSubTaskContext struct {
	*defaultExecContext
	taskCtx *DefaultTaskContext
}

func NewDefaultTaskContext(
	name string,
	ctx context.Context,
	logger core.Logger,
	data interface{},
	subtasks map[string]bool,
) core.TaskContext {
	return &DefaultTaskContext{
		newDefaultExecContext(name, ctx, data, logger),
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

var _ core.TaskContext = (*DefaultTaskContext)(nil)

func (c *DefaultSubTaskContext) TaskContext() core.TaskContext {
	return c.taskCtx
}

var _ core.SubTaskContext = (*DefaultSubTaskContext)(nil)
