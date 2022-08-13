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

package logger

import (
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"path/filepath"
	"strings"
)

var inner *logrus.Logger
var Global core.Logger

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
	inner.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	basePath := cfg.GetString("LOGGING_DIR")
	if basePath != "" {
		basePath = filepath.Join(basePath, "devlake.log")
	}
	var err error
	Global, err = NewDefaultLogger(inner)
	Global.SetStream(&core.LoggerStreamConfig{
		Path:   basePath,
		Writer: createLogStream(basePath),
	})
	if err != nil {
		panic(err)
	}
}

func createLogStream(path string) io.Writer {
	if path == "" {
		return nil
	}
	stream, err := GetFileStream(path)
	if err != nil {
		panic(err)
	}
	return stream
}
