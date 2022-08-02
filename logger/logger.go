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
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

const defaultLogFilename = "devlake"
const defaultBasePath = "./logs"

type DefaultLogger struct {
	prefix    string
	log       *logrus.Logger
	pool      map[string]core.Logger
	basePath  string
	directory string
	filename  string
}

func NewDefaultLogger(log *logrus.Logger, rootLogDir string) (core.Logger, error) {
	if rootLogDir == "" {
		rootLogDir = defaultBasePath
	}
	defaultLogger := &DefaultLogger{
		prefix:    "",
		log:       log,
		pool:      map[string]core.Logger{},
		basePath:  rootLogDir,
		filename:  defaultLogFilename,
		directory: "",
	}
	filePath := defaultLogger.getFilePath(defaultLogger.directory, defaultLogFilename)
	err := defaultLogger.createNewLogFilestream(log, filePath)
	if err != nil {
		return nil, err
	}
	return defaultLogger, nil
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

func (l *DefaultLogger) Nested(newPrefix string, config ...*core.LoggerConfig) core.Logger {
	newTotalPrefix := newPrefix
	if newPrefix != "" {
		newTotalPrefix = fmt.Sprintf("%s [%s]", l.prefix, newPrefix)
	}
	cfg := &core.LoggerConfig{}
	if len(config) == 1 {
		cfg = config[0]
	} else if len(config) > 1 {
		panic("more than one config provided")
	}
	newLogger, err := l.getLogger(newTotalPrefix, cfg)
	if err != nil {
		l.Error("error getting a new logger: %v", newLogger)
		return l
	}
	return newLogger
}

func (l *DefaultLogger) GetFsPath() string {
	return fmt.Sprintf("%s/%s/%s.log", l.basePath, l.directory, l.filename)
}

func (l *DefaultLogger) getLogger(prefix string, config *core.LoggerConfig) (core.Logger, error) {
	// if there are zero-values, inherit from current logger
	if config.Filename == "" {
		config.Filename = l.filename
	}
	if config.Directory == "" || config.InheritBase {
		config.Directory = l.directory
	}
	logFilePath := l.getFilePath(config.Directory, config.Filename)
	loggerKey := getLoggerKey(prefix, logFilePath)
	newLogger, ok := l.pool[loggerKey]
	if ok {
		return newLogger, nil
	}
	// cache miss - create a new instance
	newLogrus := logrus.New()
	newLogrus.SetLevel(l.log.Level)
	newLogrus.SetFormatter(l.log.Formatter)
	err := l.createNewLogFilestream(newLogrus, logFilePath)
	if err != nil {
		return nil, err
	}
	newLogger = &DefaultLogger{
		prefix:    prefix,
		log:       newLogrus,
		pool:      l.pool,
		basePath:  l.basePath,
		directory: config.Directory,
		filename:  config.Filename,
	}
	l.pool[loggerKey] = newLogger
	return newLogger, nil
}

func (l *DefaultLogger) createNewLogFilestream(logger *logrus.Logger, logFilePath string) error {
	err := os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	writerStd := os.Stdout
	if file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
		logger.SetOutput(io.MultiWriter(writerStd, file))
	} else {
		return err
	}
	return nil
}

func (l *DefaultLogger) getFilePath(directory string, filename string) string {
	filename = filename + ".log"
	return filepath.Join(l.basePath, directory, filename)
}

func getLoggerKey(prefix string, filename string) string {
	return fmt.Sprintf("%s-%s", filename, prefix)
}

var _ core.Logger = (*DefaultLogger)(nil)
