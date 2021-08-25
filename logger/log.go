package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	// log.SetFormatter(&logrus.JSONFormatter{})
	// TODO: setting log level with config
	log.SetLevel(logrus.DebugLevel)
}

func Debug(data interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
	}
	log.Debug(fmt.Sprintf("[%s:%d]", file, line), data)
}

func Info(data interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
	}
	log.Info(fmt.Sprintf("[%s:%d]", file, line), data)
}

func Error(data interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
	}
	log.Error(fmt.Sprintf("[%s:%d]", file, line), data)
}

func Warn(data interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
	}
	log.Warn(fmt.Sprintf("[%s:%d]", file, line), data)
}
