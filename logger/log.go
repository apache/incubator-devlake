package logger

import (
	"fmt"
	"runtime"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var inner *logrus.Logger
var Global core.Logger

func init() {
	inner = logrus.New()
	logLevel := logrus.InfoLevel
	switch config.GetConfig().GetString("LOGGING_LEVEL") {
	case "Debug":
		logLevel = logrus.DebugLevel
	case "Info":
		logLevel = logrus.InfoLevel
	case "Warn":
		logLevel = logrus.WarnLevel
	case "Error":
		logLevel = logrus.ErrorLevel
	}
	inner.SetLevel(logLevel)
	inner.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	Global = helper.NewDefaultLogger(inner, "")
}

// TODO: remove code start from this line
var (
	Black   = Color("\033[30m%s\033[0m")
	Red     = Color("\033[31m%s\033[0m")
	Green   = Color("\033[32m%s\033[0m")
	Yellow  = Color("\033[33m%s\033[0m")
	Purple  = Color("\033[34m%s\033[0m")
	Magenta = Color("\033[35m%s\033[0m")
	Teal    = Color("\033[36m%s\033[0m")
	White   = Color("\033[37m%s\033[0m")
)

// Deprecated: color and format should be handled by logrus
func Color(colorString string) func(...interface{}) string {
	if config.GetConfig().GetBool("NO_COLOR") {
		return fmt.Sprint
	}
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

// Deprecated: use core.Logger interface instead of global variable
func Log(context string, data interface{}, color func(...interface{}) string, level string, logFunction func(args ...interface{})) {
	// This operation is likely to be slow, should be avoided except for error handling
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
	}
	logFunction(color("[", level, " >>> ", context, " - ", file, ":", line, " - "), data)
}

// Deprecated: use core.Logger interface instead of global variable
func Print(context string) {
	Log(context, nil, Magenta, "DEBUG", inner.Info)
}

// Deprecated: use core.Logger interface instead of global variable
func Debug(context string, data interface{}) {
	Log(context, data, Green, "DEBUG", inner.Debug)
}

// Deprecated: use core.Logger interface instead of global variable
func Info(context string, data interface{}) {
	Log(context, data, Teal, "INFO", inner.Info)
}

// Deprecated: use core.Logger interface instead of global variable
func Error(context string, data interface{}) {
	Log(context, data, Red, "ERROR", inner.Error)
}

// Deprecated: use core.Logger interface instead of global variable
func Warn(context string, data interface{}) {
	Log(context, data, Yellow, "WARN", inner.Warn)
}
