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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"gorm.io/gorm"
)

var taskLog = logger.Global.Nested("task service")
var activityPattern = regexp.MustCompile(`task #(\d+)`)

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
func (rt *RunningTask) Add(taskId uint64, cancel context.CancelFunc) error {
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
func (rt *RunningTask) FillProgressDetailToTasks(tasks []models.Task) {
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
func (rt *RunningTask) Remove(taskId uint64) (context.CancelFunc, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if d, ok := rt.tasks[taskId]; ok {
		delete(rt.tasks, taskId)
		return d.Cancel, nil
	}
	return nil, errors.NotFound.New(fmt.Sprintf("task with id %d not found", taskId))
}

var runningTasks RunningTask

// TaskQuery FIXME ...
type TaskQuery struct {
	Status     string `form:"status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Plugin     string `form:"plugin"`
	PipelineId uint64 `form:"pipelineId" uri:"pipelineId"`
	Pending    int    `form:"pending"`
}

func init() {
	// set all previous unfinished tasks to status failed
	runningTasks.tasks = make(map[uint64]*RunningTaskData)
}

// CreateTask FIXME ...
func CreateTask(newTask *models.NewTask) (*models.Task, error) {
	b, err := json.Marshal(newTask.Options)
	if err != nil {
		return nil, err
	}
	s, err := json.Marshal(newTask.Subtasks)
	if err != nil {
		return nil, err
	}

	task := models.Task{
		Plugin:      newTask.Plugin,
		Subtasks:    s,
		Options:     b,
		Status:      models.TASK_CREATED,
		Message:     "",
		PipelineId:  newTask.PipelineId,
		PipelineRow: newTask.PipelineRow,
		PipelineCol: newTask.PipelineCol,
	}
	err = db.Save(&task).Error
	if err != nil {
		taskLog.Error(err, "save task failed")
		return nil, errors.Internal.Wrap(err, "save task failed")
	}
	return &task, nil
}

// GetTasks FIXME ...
func GetTasks(query *TaskQuery) ([]models.Task, int64, error) {
	db := db.Model(&models.Task{}).Order("id DESC")
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Plugin != "" {
		db = db.Where("plugin = ?", query.Plugin)
	}
	if query.PipelineId > 0 {
		db = db.Where("pipeline_id = ?", query.PipelineId)
	}
	if query.Pending > 0 {
		db = db.Where("finished_at is null")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if query.Page > 0 && query.PageSize > 0 {
		offset := query.PageSize * (query.Page - 1)
		db = db.Limit(query.PageSize).Offset(offset)
	}
	tasks := make([]models.Task, 0)
	err = db.Find(&tasks).Error
	if err != nil {
		return nil, count, err
	}

	runningTasks.FillProgressDetailToTasks(tasks)

	return tasks, count, nil
}

// GetTask FIXME ...
func GetTask(taskId uint64) (*models.Task, error) {
	task := &models.Task{}
	err := db.First(task, taskId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound.New("task not found", errors.AsUserMessage())
		}
		return nil, errors.Internal.Wrap(err, "error getting the task from database", errors.AsUserMessage())
	}
	return task, nil
}

// CancelTask FIXME ...
func CancelTask(taskId uint64) error {
	cancel, err := runningTasks.Remove(taskId)
	if err != nil {
		return err
	}
	cancel()
	return nil
}

func runTasksStandalone(parentLogger core.Logger, taskIds []uint64) error {
	results := make(chan error)
	for _, taskId := range taskIds {
		taskId := taskId
		go func() {
			taskLog.Info("run task in background ", taskId)
			results <- runTaskStandalone(parentLogger, taskId)
		}()
	}
	errs := make([]string, 0)
	var err error
	finished := 0
	for err = range results {
		if err != nil {
			taskLog.Error(err, "task failed")
			errs = append(errs, err.Error())
		}
		finished++
		if finished == len(taskIds) {
			close(results)
		}
	}
	if len(errs) > 0 {
		err = errors.Default.New(strings.Join(errs, "\n"))
	}
	return err
}

func runTaskStandalone(parentLog core.Logger, taskId uint64) error {
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
	progress := make(chan core.RunningProgress, 100)
	go updateTaskProgress(taskId, progress)
	err = runner.RunTask(
		ctx,
		cfg,
		parentLog,
		db,
		progress,
		taskId,
	)
	close(progress)
	return err
}

func updateTaskProgress(taskId uint64, progress chan core.RunningProgress) {
	data := runningTasks.tasks[taskId]
	if data == nil {
		return
	}
	progressDetail := data.ProgressDetail
	for p := range progress {
		runner.UpdateProgressDetail(db, log, taskId, progressDetail, &p)
	}
}

func getTaskIdFromActivityId(activityId string) (uint64, error) {
	submatches := activityPattern.FindStringSubmatch(activityId)
	if len(submatches) < 2 {
		return 0, errors.Default.New("activityId does not match")
	}
	return strconv.ParseUint(submatches[1], 10, 64)
}
