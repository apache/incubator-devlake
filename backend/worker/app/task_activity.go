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
	"context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/runner"
	"go.temporal.io/sdk/activity"
)

// DevLakeTaskActivity FIXME ...
func DevLakeTaskActivity(ctx context.Context, configJson []byte, taskId uint64, loggerConfig *log.LoggerConfig) errors.Error {
	basicRes, err := loadResources(configJson, loggerConfig)
	if err != nil {
		return err
	}
	logger := basicRes.GetLogger()
	logger.Info("received task #%d", taskId)
	progressDetail := &models.TaskProgressDetail{}
	progChan := make(chan plugin.RunningProgress)
	defer close(progChan)
	go func() {
		for p := range progChan {
			runner.UpdateProgressDetail(basicRes, taskId, progressDetail, &p)
			activity.RecordHeartbeat(ctx, progressDetail)
		}
	}()
	err = runner.RunTask(ctx, basicRes, progChan, taskId)
	if err != nil {
		logger.Error(err, "failed to execute task #%d", taskId)
	}
	logger.Info("finished task #%d", taskId)
	return err
}
