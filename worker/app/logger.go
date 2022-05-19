package app

import (
	"fmt"

	"github.com/apache/incubator-devlake/plugins/core"
	"go.temporal.io/sdk/log"
)

type TemporalLogger struct {
	log core.Logger
}

func NewTemporalLogger(log core.Logger) log.Logger {
	return &TemporalLogger{
		log,
	}
}

func (l *TemporalLogger) Log(lv core.LogLevel, msg string, keyvals ...interface{}) {
	if l.log.IsLevelEnabled(lv) {
		for i := 0; i < len(keyvals); i += 2 {
			msg += fmt.Sprintf(" %s %v", keyvals[i], keyvals[i+1])
		}
		l.log.Log(lv, msg)
	}
}

func (l *TemporalLogger) Debug(msg string, keyvals ...interface{}) {
	l.Log(core.LOG_DEBUG, msg, keyvals...)
}

func (l *TemporalLogger) Info(msg string, keyvals ...interface{}) {

	l.Log(core.LOG_INFO, msg, keyvals...)
}

func (l *TemporalLogger) Warn(msg string, keyvals ...interface{}) {

	l.Log(core.LOG_WARN, msg, keyvals...)
}

func (l *TemporalLogger) Error(msg string, keyvals ...interface{}) {
	l.Log(core.LOG_ERROR, msg, keyvals...)
}
