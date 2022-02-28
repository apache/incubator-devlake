package logger

import (
	"fmt"
	"runtime"

	"github.com/merico-dev/lake/config"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	// log.SetFormatter(&logrus.JSONFormatter{})
	// TODO: setting log level with config
	log.SetLevel(logrus.DebugLevel)
}

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

func Log(context string, data interface{}, color func(...interface{}) string, level string, logFunction func(args ...interface{})) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
	}
	logFunction(color("[", level, " >>> ", context, " - ", file, ":", line, " - "), data)
}

func Print(context string) {
	Log(context, nil, Magenta, "DEBUG", log.Info)
}

func Debug(context string, data interface{}) {
	Log(context, data, Green, "DEBUG", log.Debug)
}

func Info(context string, data interface{}) {
	Log(context, data, Teal, "INFO", log.Info)
}

func Error(context string, data interface{}) {
	Log(context, data, Red, "ERROR", log.Error)
}

func Warn(context string, data interface{}) {
	Log(context, data, Yellow, "WARN", log.Warn)
}

func GetLogger() *logrus.Logger {
	return log
}
