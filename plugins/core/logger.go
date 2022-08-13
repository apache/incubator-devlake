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

package core

import (
	"github.com/sirupsen/logrus"
	"io"
)

type LogLevel logrus.Level

const (
	LOG_DEBUG LogLevel = LogLevel(logrus.DebugLevel)
	LOG_INFO  LogLevel = LogLevel(logrus.InfoLevel)
	LOG_WARN  LogLevel = LogLevel(logrus.WarnLevel)
	LOG_ERROR LogLevel = LogLevel(logrus.ErrorLevel)
)

// Logger General logger interface, can be used anywhere
type Logger interface {
	IsLevelEnabled(level LogLevel) bool
	Printf(format string, a ...interface{})
	Log(level LogLevel, format string, a ...interface{})
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	// Nested return a new logger instance. `name` is the extra prefix to be prepended to each message. Leaving it blank
	// will add no additional prefix. The new Logger will inherit the properties of the original.
	Nested(name string) Logger
	// GetConfig Returns a copy of the LoggerConfig associated with this Logger. This is meant to be used by the framework.
	GetConfig() *LoggerConfig
	// SetStream sets the output of this Logger. This is meant to be used by the framework.
	SetStream(config *LoggerStreamConfig)
}

// LoggerStreamConfig stream related config to set on a Logger
type LoggerStreamConfig struct {
	Path   string
	Writer io.Writer
}

type InjectLogger interface {
	SetLogger(logger Logger)
}

// LoggerConfig config related to the Logger. This needs to be serializable, so it can be passed around over the wire.
type LoggerConfig struct {
	Path   string
	Prefix string
}
