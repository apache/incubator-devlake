package helper

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/sirupsen/logrus"
)

// bridge to current implementation at this point
// TODO: implement another TaskLogger for distributed runner/worker
type DefaultLogger struct {
	prefix string
	log    *logrus.Logger
}

func NewDefaultTaskLogger(log *logrus.Logger, prefix string) *DefaultLogger {
	if log == nil {
		log = logger.GetLogger()
	}
	return &DefaultLogger{prefix: prefix, log: log}

}

func (l *DefaultLogger) Log(level core.LogLevel, format string, a ...interface{}) {
	lv := logrus.Level(level)
	if l.log.IsLevelEnabled(lv) {
		msg := fmt.Sprintf(format, a...)
		if l.prefix != "" {
			msg = fmt.Sprintf("%s %s", l.prefix, msg)
		}
		l.log.Log(lv, msg)
	}
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
	return NewDefaultTaskLogger(l.log, fmt.Sprintf("%s [%s]", l.prefix, name))
}

var _ core.Logger = (*DefaultLogger)(nil)
