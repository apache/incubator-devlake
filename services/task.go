package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/runner"
)

type RunningTaskData struct {
	Cancel         context.CancelFunc
	ProgressDetail *models.TaskProgressDetail
}

type RunningTask struct {
	mu    sync.Mutex
	tasks map[uint64]*RunningTaskData
}

func taskServiceInit() {
	// reset task status
	db.Model(&models.Task{}).Where("status <> ?", models.TASK_COMPLETED).Update("status", models.TASK_FAILED)
}

func (rt *RunningTask) Add(taskId uint64, cancel context.CancelFunc) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if _, ok := rt.tasks[taskId]; ok {
		return fmt.Errorf("task with id %v already running", taskId)
	}
	rt.tasks[taskId] = &RunningTaskData{
		Cancel:         cancel,
		ProgressDetail: &models.TaskProgressDetail{},
	}
	return nil
}

// less lock times than GetProgressDetail
func (rt *RunningTask) GetProgressDetailToTasks(tasks []models.Task) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	for index, task := range tasks {
		taskId := task.ID
		if task, ok := rt.tasks[taskId]; ok {
			tasks[index].ProgressDetail = task.ProgressDetail
		}
	}

	return nil
}

func (rt *RunningTask) GetProgressDetail(taskId uint64) *models.TaskProgressDetail {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if task, ok := rt.tasks[taskId]; ok {
		return task.ProgressDetail
	}
	return nil
}

func (rt *RunningTask) Remove(taskId uint64) (context.CancelFunc, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if d, ok := rt.tasks[taskId]; ok {
		delete(rt.tasks, taskId)
		return d.Cancel, nil
	}
	return nil, fmt.Errorf("task with id %v not found", taskId)
}

var runningTasks RunningTask

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

func CreateTask(newTask *models.NewTask) (*models.Task, error) {
	b, err := json.Marshal(newTask.Options)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		Plugin:      newTask.Plugin,
		Options:     b,
		Status:      models.TASK_CREATED,
		Message:     "",
		PipelineId:  newTask.PipelineId,
		PipelineRow: newTask.PipelineRow,
		PipelineCol: newTask.PipelineCol,
	}
	err = db.Save(&task).Error
	if err != nil {
		logger.Error("save task failed", err)
		return nil, errors.InternalError
	}
	return &task, nil
}

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

	runningTasks.GetProgressDetailToTasks(tasks)

	return tasks, count, nil
}

func GetTask(taskId uint64) (*models.Task, error) {
	task := &models.Task{}
	err := db.Find(task, taskId).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func CancelTask(taskId uint64) error {
	cancel, err := runningTasks.Remove(taskId)
	if err != nil {
		return err
	}
	cancel()
	return nil
}

func runTasksStandalone(taskIds []uint64) error {
	results := make(chan error)
	for _, taskId := range taskIds {
		taskId := taskId
		go func() {
			log.Info("run task in background ", taskId)
			results <- runTaskStandalone(taskId)
		}()
	}
	errs := make([]string, 0)
	var err error
	finished := 0
	for err = range results {
		if err != nil {
			log.Error("pipeline task failed", err)
			errs = append(errs, err.Error())
		}
		if finished == len(taskIds) {
			close(results)
		}
	}
	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "\n"))
	}
	return err
}

func runTaskStandalone(taskId uint64) error {
	// deferng cleaning up
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
		cfg,
		logger.Global.Nested(fmt.Sprintf("task #%d", taskId)),
		db,
		ctx,
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
	task := &models.Task{}
	task.ID = taskId
	for p := range progress {
		switch p.Type {
		case core.TaskSetProgress:
			progressDetail.TotalSubTasks = p.Total
			progressDetail.FinishedSubTasks = p.Current
		case core.TaskIncProgress:
			progressDetail.FinishedSubTasks = p.Current
			// TODO: get rid of db update
			pct := float32(p.Current) / float32(p.Total)
			err := db.Model(task).Update("progress", pct).Error
			if err != nil {
				logger.Global.Error("failed to update progress: %w", err)
			}
		case core.SubTaskSetProgress:
			progressDetail.TotalRecords = p.Total
			progressDetail.FinishedRecords = p.Current
		case core.SubTaskIncProgress:
			progressDetail.FinishedRecords = p.Current
		case core.SetCurrentSubTask:
			progressDetail.SubTaskName = p.SubTaskName
			progressDetail.SubTaskNumber = p.SubTaskNumber
		}
	}
}
