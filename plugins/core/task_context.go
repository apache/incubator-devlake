package core

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// prepare for temporal

type LogLevel logrus.Level

const (
	LOG_DEBUG LogLevel = LogLevel(logrus.DebugLevel)
	LOG_INFO  LogLevel = LogLevel(logrus.InfoLevel)
	LOG_WARN  LogLevel = LogLevel(logrus.WarnLevel)
	LOG_ERROR LogLevel = LogLevel(logrus.ErrorLevel)
)

type Logger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
}

type TaskLogger interface {
	Logger
	// update progress, pass -1 for total if it was unavailable
	Progress(subtask string, current int, total int)
}

type DefaultTaskLogger struct {
	logger *logrus.Logger
}

func (l *DefaultTaskLogger) Log(level LogLevel, format string, a ...interface{}) {
	lv := logrus.Level(level)
	if l.logger.IsLevelEnabled(lv) {
		l.logger.Log(lv, fmt.Sprintf(format, a...))
	}
}

func (l *DefaultTaskLogger) Debug(format string, a ...interface{}) {
	l.Log(LOG_DEBUG, format, a...)
}

func (l *DefaultTaskLogger) Info(format string, a ...interface{}) {
	l.Log(LOG_INFO, format, a...)
}

func (l *DefaultTaskLogger) Warn(format string, a ...interface{}) {
	l.Log(LOG_WARN, format, a...)
}

func (l *DefaultTaskLogger) Error(format string, a ...interface{}) {
	l.Log(LOG_ERROR, format, a...)
}

func (l *DefaultTaskLogger) Progress(subtask string, current int, total int) {
	l.Info("progress of subtask %v is updated: %v done / %v total", subtask, current, total)
}

type TaskContext interface {
	GetConfig(name string) string
	GetDb() *gorm.DB
	GetContext() context.Context
	GetData() interface{}
	GetLogger() TaskLogger
}

// bridge to current implementation at this point
type DefaultTaskContext struct {
	ctx    context.Context
	data   interface{}
	logger TaskLogger
}

func NewDefaultTaskContext(ctx context.Context, data interface{}) TaskContext {
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

func (c *DefaultTaskContext) GetLogger() TaskLogger {
	return c.logger
}

// update progress, pass -1 for total if it was unavailable
func (c *DefaultTaskContext) Progress(subtask string, current int, total int) {
	panic("not implemented") // TODO: Implement
}

var _ TaskContext = (*DefaultTaskContext)(nil)
