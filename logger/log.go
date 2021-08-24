package logger

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func init() {
	log = logrus.New()
	// log.SetFormatter(&logrus.JSONFormatter{})
	// TODO: setting log level with config
	log.SetLevel(logrus.DebugLevel)
}

func Debug(message string, ctx interface{}) {
	log.Debug(message, ctx)
}

func Info(message string, ctx interface{}) {
	log.Info(message, ctx)
}

func Error(message string, err error) {
	log.Error(message, err)
}

func Warn(message string, ctx interface{}) {
	log.Warn(message, ctx)
}
