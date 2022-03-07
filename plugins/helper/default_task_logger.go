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
	log *logrus.Logger
}

func NewDefaultTaskLogger(log *logrus.Logger) *DefaultLogger {
	if log == nil {
		log = logger.GetLogger()
	}
	return &DefaultLogger{log: log}

}

func (l *DefaultLogger) Log(level core.LogLevel, format string, a ...interface{}) {
	lv := logrus.Level(level)
	if l.log.IsLevelEnabled(lv) {
		l.log.Log(lv, fmt.Sprintf(format, a...))
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

var _ core.Logger = (*DefaultLogger)(nil)
