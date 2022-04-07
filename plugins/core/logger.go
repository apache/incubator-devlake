package core

import "github.com/sirupsen/logrus"

type LogLevel logrus.Level

const (
	LOG_DEBUG LogLevel = LogLevel(logrus.DebugLevel)
	LOG_INFO  LogLevel = LogLevel(logrus.InfoLevel)
	LOG_WARN  LogLevel = LogLevel(logrus.WarnLevel)
	LOG_ERROR LogLevel = LogLevel(logrus.ErrorLevel)
)

// General logger interface, can be used any where
type Logger interface {
	IsLevelEnabled(level LogLevel) bool
	Printf(format string, a ...interface{})
	Log(level LogLevel, format string, a ...interface{})
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	// return a new logger which output nested log
	Nested(name string) Logger
}
