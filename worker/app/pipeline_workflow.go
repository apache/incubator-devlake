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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/runner"
	"go.temporal.io/sdk/workflow"
)

// DevLakePipelineWorkflow FIXME ...
func DevLakePipelineWorkflow(ctx workflow.Context, configJson []byte, pipelineId uint64) error {
	cfg, logger, db, err := loadResources(configJson)
	if err != nil {
		return err
	}
	logger.Info("received pipeline #%d", pipelineId)
	err = runner.RunPipeline(
		cfg,
		logger,
		db,
		pipelineId,
		func(taskIds []uint64) error {
			futures := make([]workflow.Future, len(taskIds))
			for i, taskId := range taskIds {
				activityOpts := workflow.ActivityOptions{
					ActivityID:          fmt.Sprintf("task #%d", taskId),
					StartToCloseTimeout: 24 * time.Hour,
					WaitForCancellation: true,
				}
				activityCtx := workflow.WithActivityOptions(ctx, activityOpts)
				futures[i] = workflow.ExecuteActivity(activityCtx, DevLakeTaskActivity, configJson, taskId)
			}
			errs := make([]string, 0)
			for _, future := range futures {
				err := future.Get(ctx, nil)
				if err != nil {
					errs = append(errs, err.Error())
				}
			}
			if len(errs) > 0 {
				return fmt.Errorf(strings.Join(errs, "\n"))
			}
			return nil
		},
	)
	if err != nil {
		logger.Error("failed to execute pipeline #%d: %w", pipelineId, err)
	}
	logger.Info("finished pipeline #%d", pipelineId)
	return err
}
