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
	"regexp"
	"strings"
)

var alreadyInBracketsRegex = regexp.MustCompile(`\[.*?]+`)

type DefaultLogger struct {
	log    *logrus.Logger
	config *core.LoggerConfig
}

func NewDefaultLogger(log *logrus.Logger) (core.Logger, error) {
	defaultLogger := &DefaultLogger{
		log:    log,
		config: &core.LoggerConfig{},
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
		if l.config.Prefix != "" {
			msg = fmt.Sprintf("%s %s", l.config.Prefix, msg)
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

func (l *DefaultLogger) SetStream(config *core.LoggerStreamConfig) {
	if config.Path != "" {
		l.config.Path = config.Path
	}
	if config.Writer != nil {
		l.log.SetOutput(config.Writer)
	}
}

func (l *DefaultLogger) GetConfig() *core.LoggerConfig {
	return &core.LoggerConfig{
		Path:   l.config.Path,
		Prefix: l.config.Prefix,
	}
}

func (l *DefaultLogger) Nested(newPrefix string) core.Logger {
	newTotalPrefix := newPrefix
	if newPrefix != "" {
		newTotalPrefix = l.createPrefix(newPrefix)
	}
	newLogger, err := l.getLogger(newTotalPrefix)
	if err != nil {
		l.Error("error getting a new logger: %v", newLogger)
		return l
	}
	return newLogger
}

func (l *DefaultLogger) getLogger(prefix string) (core.Logger, error) {
	newLogrus := logrus.New()
	newLogrus.SetLevel(l.log.Level)
	newLogrus.SetFormatter(l.log.Formatter)
	newLogrus.SetOutput(l.log.Out)
	newLogger := &DefaultLogger{
		log: newLogrus,
		config: &core.LoggerConfig{
			Path:   l.config.Path,
			Prefix: prefix,
		},
	}
	return newLogger, nil
}

func (l *DefaultLogger) createPrefix(newPrefix string) string {
	newPrefix = strings.TrimSpace(newPrefix)
	alreadyInBrackets := alreadyInBracketsRegex.MatchString(newPrefix)
	if alreadyInBrackets {
		return fmt.Sprintf("%s %s", l.config.Prefix, newPrefix)
	}
	return fmt.Sprintf("%s [%s]", l.config.Prefix, newPrefix)
}

var _ core.Logger = (*DefaultLogger)(nil)
