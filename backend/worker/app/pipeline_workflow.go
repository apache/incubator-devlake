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
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/runner"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"
)

// DevLakePipelineWorkflow FIXME ...
func DevLakePipelineWorkflow(ctx workflow.Context, configJson []byte, pipelineId uint64, loggerConfig *log.LoggerConfig) errors.Error {
	basicRes, err := loadResources(configJson, loggerConfig)
	if err != nil {
		return errors.Convert(err)
	}
	logger := basicRes.GetLogger()
	logger.Info("received pipeline #%d", pipelineId)
	err = runner.RunPipeline(
		basicRes,
		pipelineId,
		func(taskIds []uint64) errors.Error {
			return runTasks(ctx, configJson, taskIds, logger)
		},
	)
	if err != nil {
		logger.Error(err, "failed to execute pipeline #%d", pipelineId)
	}
	logger.Info("finished pipeline #%d", pipelineId)
	return err
}

func runTasks(ctx workflow.Context, configJson []byte, taskIds []uint64, logger log.Logger) errors.Error {
	cleanExit := false
	defer func() {
		if !cleanExit {
			logger.Error(nil, "fatal error while executing task Ids: %v", taskIds)
		}
	}()
	futures := make([]workflow.Future, len(taskIds))
	for i, taskId := range taskIds {
		activityOpts := workflow.ActivityOptions{
			ActivityID:          fmt.Sprintf("task #%d", taskId),
			StartToCloseTimeout: 24 * time.Hour,
			WaitForCancellation: true,
		}
		activityCtx := workflow.WithActivityOptions(ctx, activityOpts)
		futures[i] = workflow.ExecuteActivity(activityCtx, DevLakeTaskActivity, configJson, taskId, logger)
	}
	errs := make([]string, 0)
	for _, future := range futures {
		err := future.Get(ctx, nil)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	cleanExit = true
	if len(errs) > 0 {
		return errors.Default.New(strings.Join(errs, "\n"))
	}
	return nil
}
