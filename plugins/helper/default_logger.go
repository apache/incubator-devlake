package helper

import (
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
)

type DefaultLogger struct {
	prefix     string
	log        *logrus.Logger
	loggerPool map[string]*logrus.Logger
}

func NewDefaultLogger(log *logrus.Logger, prefix string, loggerPool map[string]*logrus.Logger) *DefaultLogger {
	newDefaultLogger := &DefaultLogger{prefix: prefix, log: log}
	newDefaultLogger.loggerPool = loggerPool
	return newDefaultLogger
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

// bind two writer to logger
func (l *DefaultLogger) Nested(name string) core.Logger {
	writerStd := os.Stdout
	fileName := ""
	loggerPrefixRegex := regexp.MustCompile(`(\[task #\d+]\s\[\w+])`)
	groups := loggerPrefixRegex.FindStringSubmatch(fmt.Sprintf("%s [%s]", l.prefix, name))
	if len(groups) > 1 {
		fileName = groups[1]
	}

	if fileName == "" {
		fileName = "devlake"
	}

	if l.loggerPool[fileName] != nil {
		return NewDefaultLogger(l.loggerPool[fileName], fmt.Sprintf("%s [%s]", l.prefix, name), l.loggerPool)
	}
	newLog := logrus.New()
	newLog.SetLevel(l.log.Level)
	newLog.SetFormatter(l.log.Formatter)

	if file, err := os.OpenFile(fmt.Sprintf("logs/%s.log", fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
		newLog.SetOutput(io.MultiWriter(writerStd, file))
	}
	l.loggerPool[fileName] = newLog
	return NewDefaultLogger(newLog, fmt.Sprintf("%s [%s]", l.prefix, name), l.loggerPool)
}

var _ core.Logger = (*DefaultLogger)(nil)
