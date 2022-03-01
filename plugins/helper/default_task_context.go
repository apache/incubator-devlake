package helper

import (
	"context"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
)

// bridge to current implementation at this point
type DefaultTaskContext struct {
	ctx    context.Context
	data   interface{}
	logger core.TaskLogger
}

func NewDefaultTaskContext(ctx context.Context, data interface{}) core.TaskContext {
	return &DefaultTaskContext{
		ctx,
		data,
		&DefaultTaskLogger{},
	}
}

func (c *DefaultTaskContext) GetConfig(name string) string {
	return config.GetConfig().GetString(name)
}

func (c *DefaultTaskContext) GetDb() *gorm.DB {
	return models.Db
}

func (c *DefaultTaskContext) GetContext() context.Context {
	return c.ctx
}

func (c *DefaultTaskContext) GetData() interface{} {
	return c.data
}

func (c *DefaultTaskContext) GetLogger() core.TaskLogger {
	return c.logger
}

// update progress, pass -1 for total if it was unavailable
func (c *DefaultTaskContext) Progress(subtask string, current int, total int) {
	panic("not implemented") // TODO: Implement
}

var _ core.TaskContext = (*DefaultTaskContext)(nil)
