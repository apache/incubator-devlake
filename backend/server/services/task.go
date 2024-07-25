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
	"math"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/impls/logruslog"
)

var taskLog = logruslog.Global.Nested("task service")

// TaskQuery FIXME .
type TaskQuery struct {
	Pagination
	Status     string `form:"status"`
	Plugin     string `form:"plugin"`
	PipelineId uint64 `form:"pipelineId" uri:"pipelineId"`
	Pending    int    `form:"pending"`
}

func createTask(newTask *models.NewTask, tx dal.Transaction) (*models.Task, errors.Error) {
	task := &models.Task{
		Plugin:      newTask.Plugin,
		Subtasks:    newTask.Subtasks,
		Options:     newTask.Options,
		Status:      models.TASK_CREATED,
		Message:     "",
		PipelineId:  newTask.PipelineId,
		PipelineRow: newTask.PipelineRow,
		PipelineCol: newTask.PipelineCol,
	}
	if newTask.IsRerun {
		task.Status = models.TASK_RERUN
	}
	err := tx.Create(task)
	if err != nil {
		taskLog.Error(err, "save task failed")
		return nil, errors.Internal.Wrap(err, "save task failed")
	}
	return task, nil
}

// GetTasks returns paginated tasks that match the given query
func GetTasks(query *TaskQuery) ([]*models.Task, int64, errors.Error) {
	// verify query
	if err := VerifyStruct(query); err != nil {
		return nil, 0, err
	}

	// construct common query clauses
	clauses := []dal.Clause{dal.From(&models.Task{})}
	if query.Status != "" {
		clauses = append(clauses, dal.Where("status = ?", query.Status))
	}
	if query.Plugin != "" {
		clauses = append(clauses, dal.Where("plugin = ?", query.Plugin))
	}
	if query.PipelineId > 0 {
		clauses = append(clauses, dal.Where("pipeline_id = ?", query.PipelineId))
	}
	if query.Pending > 0 {
		clauses = append(clauses, dal.Where("finished_at is null"))
	}

	// count total records
	count, err := db.Count(clauses...)
	if err != nil {
		return nil, 0, err
	}

	// load paginated records from db
	clauses = append(clauses,
		dal.Orderby("id DESC"),
		dal.Offset(query.GetSkip()),
		dal.Limit(query.GetPageSizeOr(10000)),
	)
	tasks := make([]*models.Task, 0)
	err = db.All(&tasks, clauses...)
	if err != nil {
		return nil, count, err
	}

	// fill running information
	runningTasks.FillProgressDetailToTasks(tasks)

	return tasks, count, nil
}

// GetTasksWithLastStatus returns task list of the pipeline, only the most recently tasks would be returned
// TODO: adopts GetLatestTasksOfPipeline
func GetTasksWithLastStatus(pipelineId uint64, shouldSanitize bool, tx dal.Dal) ([]*models.Task, errors.Error) {
	if tx == nil {
		tx = db
	}
	var tasks []*models.Task
	err := tx.All(&tasks, dal.Where("pipeline_id = ?", pipelineId), dal.Orderby("id DESC"))
	if err != nil {
		return nil, err
	}
	taskIds := make(map[int64]struct{})
	var result []*models.Task
	var maxRow, maxCol int
	for _, task := range tasks {
		if task.PipelineRow > maxRow {
			maxRow = task.PipelineRow
		}
		if task.PipelineCol > maxCol {
			maxCol = task.PipelineCol
		}
	}

	for _, task := range tasks {
		index := int64(task.PipelineRow)*int64(maxCol) + int64(task.PipelineCol)
		if shouldSanitize {
			taskOption, err := SanitizePluginOption(task.Plugin, task.Options)
			if err != nil {
				return nil, errors.Convert(err)
			}
			task.Options = taskOption
		}
		if _, ok := taskIds[index]; !ok {
			taskIds[index] = struct{}{}
			result = append(result, task)
		}
	}

	runningTasks.FillProgressDetailToTasks(result)
	return result, nil
}

// GetTask FIXME ...
func GetTask(taskId uint64) (*models.Task, errors.Error) {
	task := &models.Task{}
	err := db.First(task, dal.Where("id = ?", taskId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.New("task not found")
		}
		return nil, errors.Internal.Wrap(err, "error getting the task from database")
	}
	return task, nil
}

// CancelTask FIXME ...
func CancelTask(taskId uint64) errors.Error {
	cancel, err := runningTasks.Remove(taskId)
	if err != nil {
		return err
	}
	cancel()
	return nil
}

// RunTasksStandalone run tasks in parallel
func RunTasksStandalone(parentLogger log.Logger, taskIds []uint64) errors.Error {
	if len(taskIds) == 0 {
		return nil
	}
	results := make(chan error)
	for _, taskId := range taskIds {
		go func(id uint64) {
			taskLog.Info("run task #%d in background ", id)
			var err errors.Error
			taskErr := runTaskStandalone(parentLogger, id)
			if taskErr != nil {
				err = errors.Default.Wrap(taskErr, fmt.Sprintf("Error running task %d.", id))
			}
			results <- err
		}(taskId)
	}
	errs := make([]error, 0)
	var err error
	finished := 0
	for err = range results {
		if err != nil {
			taskLog.Error(err, "task failed")
			errs = append(errs, err)
		}
		finished++
		if finished == len(taskIds) {
			close(results)
		}
	}
	if len(errs) > 0 {
		var sb strings.Builder
		for _, e := range errs {
			_, _ = sb.WriteString(e.Error())
			_, _ = sb.WriteString("\n")
			if errors.Is(e, context.Canceled) {
				parentLogger.Info("task canceled")
				return errors.Convert(e)
			}
		}
		err = errors.Default.New(sb.String())
	}
	return errors.Convert(err)
}

// RerunTask reruns specified task
func RerunTask(taskId uint64) (*models.Task, errors.Error) {
	task, err := GetTask(taskId)
	if err != nil {
		return nil, err
	}
	rerunTasks, err := RerunPipeline(task.PipelineId, task)
	if err != nil {
		return nil, err
	}
	rerunTask := rerunTasks[0]
	taskOption, sanitizePluginOptionErr := SanitizePluginOption(rerunTask.Plugin, rerunTask.Options)
	if sanitizePluginOptionErr != nil {
		return nil, errors.Convert(err)
	}
	rerunTask.Options = taskOption
	return rerunTask, nil
}

// GetSubTasksInfo returns subtask list of the pipeline, only the most recently subtasks would be returned
func GetSubTasksInfo(pipelineId uint64, shouldSanitize bool, tx dal.Dal) (*models.SubTasksOuput, errors.Error) {
	if tx == nil {
		tx = db
	}
	var tasks []*models.Task
	err := tx.All(&tasks, dal.Where("pipeline_id = ?", pipelineId), dal.Orderby("id"))
	if err != nil {
		return nil, err
	}

	lastTasks := filterTasksWithLastStatus(tasks)
	var subtasksInfo []models.SubtasksInfo
	var totalSubtasksCount int64
	var totalFinishedSubTasksCount int64
	var count int64
	var status []string
	for _, task := range lastTasks {
		if task.Plugin == "org" || task.Plugin == "refdiff" || task.Plugin == "dora" {
			continue
		}
		subTaskResult := models.SubtasksInfo{
			ID:           task.ID,
			PipelineID:   task.PipelineId,
			CreatedAt:    task.CreatedAt,
			UpdatedAt:    task.UpdatedAt,
			BeganAt:      task.BeganAt,
			FinishedAt:   task.FinishedAt,
			Plugin:       task.Plugin,
			Status:       task.Status,
			Message:      task.Message,
			ErrorName:    task.ErrorName,
			SpentSeconds: task.SpentSeconds,
		}
		if shouldSanitize {
			taskOption, err := SanitizePluginOption(task.Plugin, task.Options)
			if err != nil {
				return nil, errors.Convert(err)
			}
			subTaskResult.Options = taskOption
		}

		var totalTransform int64
		var finishedTransform int64
		subtasks := []*models.Subtask{}
		err = tx.All(&subtasks, dal.Where("task_id = ?", task.ID))
		if err != nil {
			return nil, err
		}
		for _, subtask := range subtasks {
			t := &models.SubtaskDetails{
				ID:              subtask.ID,
				CreatedAt:       subtask.CreatedAt,
				UpdatedAt:       subtask.UpdatedAt,
				TaskID:          subtask.TaskID,
				Name:            subtask.Name,
				Number:          subtask.Number,
				BeganAt:         subtask.BeganAt,
				FinishedAt:      subtask.FinishedAt,
				SpentSeconds:    subtask.SpentSeconds,
				FinishedRecords: subtask.FinishedRecords,
				Sequence:        subtask.Sequence,
				IsCollector:     subtask.IsCollector,
				IsFailed:        subtask.IsFailed,
				Message:         subtask.Message,
			}
			if !subtask.IsCollector {
				totalTransform++
				if subtask.FinishedAt != nil {
					finishedTransform++
				}
			}

			subTaskResult.SubtaskDetails = append(subTaskResult.SubtaskDetails, t)
		}
		subTaskResult.TotalTransform = totalTransform
		subTaskResult.FinishedTransform = finishedTransform
		subtasksInfo = append(subtasksInfo, subTaskResult)

		collectSubtasksCount := errors.Must1(tx.Count(dal.From("_devlake_subtasks"), dal.Where("task_id = ? and is_collector = true", task.ID)))
		totalSubtasksCount += collectSubtasksCount
		finishedSubTasksCount := errors.Must1(tx.Count(dal.From("_devlake_subtasks"), dal.Where("task_id = ? and is_collector = true and finished_at is not null", task.ID)))
		totalFinishedSubTasksCount += finishedSubTasksCount
		count++

		status = append(status, task.Status)
	}

	subTasksOuput := &models.SubTasksOuput{}

	subTasksOuput.SubtasksInfo = subtasksInfo
	subTasksOuput.Count = totalSubtasksCount + 1 // +1: add transform task
	completionRateFloat := float64(totalFinishedSubTasksCount) / float64(totalSubtasksCount)
	roundedCompletionRate := math.Round(completionRateFloat*100) / 100
	if math.IsNaN(roundedCompletionRate) {
		roundedCompletionRate = 1
	}
	subTasksOuput.CompletionRate = roundedCompletionRate

	subTasksOuput.Status = getTaskStatus(status)
	subTasksOuput.Count = count

	return subTasksOuput, nil
}

// filterTasksWithLastStatus returns the latest task for each plugin
func filterTasksWithLastStatus(tasks []*models.Task) []*models.Task {
	taskMap := make(map[string]*models.Task)
	for _, task := range tasks {
		key := fmt.Sprintf("%d-%d-%d", task.PipelineId, task.PipelineRow, task.PipelineCol)
		if existingTask, ok := taskMap[key]; ok {
			if task.BeganAt != nil && (existingTask.BeganAt == nil || task.BeganAt.After(*existingTask.BeganAt)) {
				taskMap[key] = task
			}
		} else {
			taskMap[key] = task
		}
	}

	var filteredTasks []*models.Task
	for _, task := range tasks {
		key := fmt.Sprintf("%d-%d-%d", task.PipelineId, task.PipelineRow, task.PipelineCol)
		if filteredTask, ok := taskMap[key]; ok {
			filteredTasks = append(filteredTasks, filteredTask)
		}
	}

	return filteredTasks
}

func getTaskStatus(statuses []string) string {
	var status string
	if len(statuses) == 0 {
		return status
	}

	failedCount := 0
	completedCount := 0
	for _, s := range statuses {
		if s == models.TASK_FAILED {
			failedCount++
		} else if s == models.TASK_COMPLETED {
			completedCount++
		}
	}
	if failedCount > 0 && completedCount > 0 {
		status = "TASK_PARTIAL"
	} else if failedCount == len(statuses) {
		status = models.TASK_FAILED
	} else if completedCount == len(statuses) {
		status = models.TASK_COMPLETED
	} else {
		status = models.TASK_RUNNING
	}

	return status
}
