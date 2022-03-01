package core

import (
	"context"

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

type TaskContext interface {
	GetConfig(name string) string
	GetDb() *gorm.DB
	GetContext() context.Context
	GetData() interface{}
	GetLogger() TaskLogger
}
