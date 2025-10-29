/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logruslog

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/sirupsen/logrus"
)

var inner *logrus.Logger
var Global log.Logger

func init() {
	inner = logrus.New()
	logLevel := logrus.InfoLevel
	cfg := config.GetConfig()
	switch strings.ToLower(cfg.GetString("LOGGING_LEVEL")) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	}
	inner.SetLevel(logLevel)

	var formatter logrus.Formatter

	format := os.Getenv("LOGGING_FORMAT")

	switch format {
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.DateTime,
		}
	default:
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			FullTimestamp:   true,
		}
	}

	inner.SetFormatter(formatter)

	basePath := cfg.GetString("LOGGING_DIR")
	if basePath == "" {
		basePath = "./logs"
	}
	abs, absErr := filepath.Abs(basePath)
	if absErr != nil {
		panic(absErr)
	}
	basePath = filepath.Join(abs, "devlake.log")
	var err errors.Error
	Global, err = NewDefaultLogger(inner)
	if err != nil {
		panic(err)
	}
	stream, err := GetFileStream(basePath)
	if err != nil {
		stream = os.Stdout
	}
	Global.SetStream(&log.LoggerStreamConfig{
		Path:   basePath,
		Writer: stream,
	})
	if err != nil {
		Global.Warn(err, "Failed to create filestream for logs. Logs will not be piped to files.")
	}
}
