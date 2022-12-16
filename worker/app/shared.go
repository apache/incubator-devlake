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

package app

import (
	"bytes"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/viper"
)

func loadResources(configJson []byte, loggerConfig *core.LoggerConfig) (core.BasicRes, errors.Error) {
	// TODO: should be redirected to server
	globalLogger := logger.Global.Nested("worker")
	// prepare
	cfg := viper.New()
	cfg.SetConfigType("json")
	err := cfg.ReadConfig(bytes.NewBuffer(configJson))
	if err != nil {
		globalLogger.Error(err, "failed to load resources")
		return nil, errors.Convert(err)
	}
	db, err := runner.NewGormDb(cfg, globalLogger)
	if err != nil {
		return nil, errors.Convert(err)
	}
	log, err := getWorkerLogger(globalLogger, loggerConfig)
	if err != nil {
		return nil, errors.Convert(err)
	}
	return runner.CreateBasicRes(cfg, log, db), nil
}

func getWorkerLogger(log core.Logger, logConfig *core.LoggerConfig) (core.Logger, errors.Error) {
	newLog := log.Nested(logConfig.Prefix)
	stream, err := logger.GetFileStream(logConfig.Path)
	if err != nil {
		return nil, err
	}
	newLog.SetStream(&core.LoggerStreamConfig{
		Path:   logConfig.Path,
		Writer: stream,
	})
	return newLog, nil
}
