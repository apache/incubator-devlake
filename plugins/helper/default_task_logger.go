package helper

import (
	"fmt"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/sirupsen/logrus"
)

// bridge to current implementation at this point
// TODO: implement another TaskLogger for distributed runner/worker
type DefaultTaskLogger struct {
	logger *logrus.Logger
}

func (l *DefaultTaskLogger) Log(level core.LogLevel, format string, a ...interface{}) {
	lv := logrus.Level(level)
	if l.logger.IsLevelEnabled(lv) {
		l.logger.Log(lv, fmt.Sprintf(format, a...))
	}
}

func (l *DefaultTaskLogger) Debug(format string, a ...interface{}) {
	l.Log(core.LOG_DEBUG, format, a...)
}

func (l *DefaultTaskLogger) Info(format string, a ...interface{}) {
	l.Log(core.LOG_INFO, format, a...)
}

func (l *DefaultTaskLogger) Warn(format string, a ...interface{}) {
	l.Log(core.LOG_WARN, format, a...)
}

func (l *DefaultTaskLogger) Error(format string, a ...interface{}) {
	l.Log(core.LOG_ERROR, format, a...)
}

func (l *DefaultTaskLogger) Progress(subtask string, current int, total int) {
	l.Info("progress of subtask %v is updated: %v done / %v total", subtask, current, total)
}

var _ core.TaskLogger = (*DefaultTaskLogger)(nil)
