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

package services

import (
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/runner"
	"strconv"
	"sync"
)

// RunningTaskData FIXME ...
type RunningTaskData struct {
	Cancel         context.CancelFunc
	ProgressDetail *models.TaskProgressDetail
}

// RunningTask FIXME ...
type RunningTask struct {
	mu    sync.Mutex
	tasks map[uint64]*RunningTaskData
}

// Add FIXME ...
func (rt *RunningTask) Add(taskId uint64, cancel context.CancelFunc) errors.Error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if _, ok := rt.tasks[taskId]; ok {
		return errors.Default.New(fmt.Sprintf("task with id %d already running", taskId))
	}
	rt.tasks[taskId] = &RunningTaskData{
		Cancel:         cancel,
		ProgressDetail: &models.TaskProgressDetail{},
	}
	return nil
}

func (rt *RunningTask) setAll(progressDetails map[uint64]*models.TaskProgressDetail) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	// delete finished tasks
	for taskId := range rt.tasks {
		if _, ok := progressDetails[taskId]; !ok {
			delete(rt.tasks, taskId)
		}
	}
	for taskId, progressDetail := range progressDetails {
		if _, ok := rt.tasks[taskId]; !ok {
			rt.tasks[taskId] = &RunningTaskData{}
		}
		rt.tasks[taskId].ProgressDetail = progressDetail
	}
}

// FillProgressDetailToTasks lock less times than GetProgressDetail
func (rt *RunningTask) FillProgressDetailToTasks(tasks []*models.Task) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	for index, task := range tasks {
		taskId := task.ID
		if task, ok := rt.tasks[taskId]; ok {
			tasks[index].ProgressDetail = task.ProgressDetail
		}
	}
}

// GetProgressDetail FIXME ...
func (rt *RunningTask) GetProgressDetail(taskId uint64) *models.TaskProgressDetail {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if task, ok := rt.tasks[taskId]; ok {
		return task.ProgressDetail
	}
	return nil
}

// Remove FIXME ...
func (rt *RunningTask) Remove(taskId uint64) (context.CancelFunc, errors.Error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if d, ok := rt.tasks[taskId]; ok {
		delete(rt.tasks, taskId)
		return d.Cancel, nil
	}
	return nil, errors.NotFound.New(fmt.Sprintf("task with id %d not found", taskId))
}

var runningTasks RunningTask

func init() {
	// set all previous unfinished tasks to status failed
	runningTasks.tasks = make(map[uint64]*RunningTaskData)
}

func runTaskStandalone(parentLog log.Logger, taskId uint64) errors.Error {
	// deferring cleaning up
	defer func() {
		_, _ = runningTasks.Remove(taskId)
	}()
	// for task cancelling
	ctx, cancel := context.WithCancel(context.Background())
	err := runningTasks.Add(taskId, cancel)
	if err != nil {
		return err
	}
	// now , create a progress update channel and kick off
	progress := make(chan plugin.RunningProgress, 100)
	go updateTaskProgress(taskId, progress)
	err = runner.RunTask(
		ctx,
		basicRes.ReplaceLogger(parentLog),
		progress,
		taskId,
	)
	close(progress)
	return err
}

func getRunningTaskById(taskId uint64) *RunningTaskData {
	runningTasks.mu.Lock()
	defer runningTasks.mu.Unlock()

	return runningTasks.tasks[taskId]
}

func updateTaskProgress(taskId uint64, progress chan plugin.RunningProgress) {
	data := getRunningTaskById(taskId)
	if data == nil {
		return
	}
	progressDetail := data.ProgressDetail
	for p := range progress {
		runningTasks.mu.Lock()
		runner.UpdateProgressDetail(basicRes, taskId, progressDetail, &p)
		runningTasks.mu.Unlock()
	}
}

func getTaskIdFromActivityId(activityId string) (uint64, errors.Error) {
	submatches := activityPattern.FindStringSubmatch(activityId)
	if len(submatches) < 2 {
		return 0, errors.Default.New("activityId does not match")
	}
	return errors.Convert01(strconv.ParseUint(submatches[1], 10, 64))
}
