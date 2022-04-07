package logger

import (
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
