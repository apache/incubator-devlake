package helper

import (
	"fmt"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/sirupsen/logrus"
)

type DefaultLogger struct {
	prefix string
	log    *logrus.Logger
}

func NewDefaultLogger(log *logrus.Logger, prefix string) *DefaultLogger {
	return &DefaultLogger{prefix: prefix, log: log}
}

func (l *DefaultLogger) IsLevelEnabled(level core.LogLevel) bool {
	if l.log == nil {
		return false
	}
	return l.log.IsLevelEnabled(logrus.Level(level))
}

func (l *DefaultLogger) Log(level core.LogLevel, format string, a ...interface{}) {
	if l.IsLevelEnabled(level) {
		msg := fmt.Sprintf(format, a...)
		if l.prefix != "" {
			msg = fmt.Sprintf("%s %s", l.prefix, msg)
		}
		l.log.Log(logrus.Level(level), msg)
	}
}

func (l *DefaultLogger) Printf(format string, a ...interface{}) {
	l.Log(core.LOG_INFO, format, a...)
}

func (l *DefaultLogger) Debug(format string, a ...interface{}) {
	l.Log(core.LOG_DEBUG, format, a...)
}

func (l *DefaultLogger) Info(format string, a ...interface{}) {
	l.Log(core.LOG_INFO, format, a...)
}

func (l *DefaultLogger) Warn(format string, a ...interface{}) {
	l.Log(core.LOG_WARN, format, a...)
}

func (l *DefaultLogger) Error(format string, a ...interface{}) {
	l.Log(core.LOG_ERROR, format, a...)
}

func (l *DefaultLogger) Nested(name string) core.Logger {
	return NewDefaultLogger(l.log, fmt.Sprintf("%s [%s]", l.prefix, name))
}

var _ core.Logger = (*DefaultLogger)(nil)
