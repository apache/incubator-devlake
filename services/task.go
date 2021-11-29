package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins"
)

type RunningTask struct {
	mu    sync.Mutex
	tasks map[uint64]context.CancelFunc
}

func (rt *RunningTask) Add(taskId uint64, cancel context.CancelFunc) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if _, ok := rt.tasks[taskId]; ok {
		return fmt.Errorf("task with id %v already running", taskId)
	}
	rt.tasks[taskId] = cancel
	return nil
}

func (rt *RunningTask) Remove(taskId uint64) (context.CancelFunc, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if cancel, ok := rt.tasks[taskId]; ok {
		delete(rt.tasks, taskId)
		return cancel, nil
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
	models.Db.Model(&models.Task{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)
	runningTasks.tasks = make(map[uint64]context.CancelFunc)
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
	err = models.Db.Save(&task).Error
	if err != nil {
		logger.Error("save task failed", err)
		return nil, errors.InternalError
	}
	return &task, nil
}

func GetTasks(query *TaskQuery) ([]models.Task, int64, error) {
	db := models.Db.Model(&models.Task{}).Order("id DESC")
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
	return tasks, count, nil
}

func GetTask(taskId uint64) (*models.Task, error) {
	task := &models.Task{}
	err := models.Db.Find(task, taskId).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

// RunTask guarantees database is update even if it panicked, and the error will be returned to caller
func RunTask(taskId uint64) error {
	task, err := GetTask(taskId)
	if err != nil {
		return err
	}
	if task.Status != models.TASK_CREATED {
		return fmt.Errorf("invalid task status")
	}

	// for task cancelling
	ctx, cancel := context.WithCancel(context.Background())
	err = runningTasks.Add(taskId, cancel)
	if err != nil {
		return err
	}

	progress := make(chan float32)
	var options map[string]interface{}
	err = json.Unmarshal(task.Options, &options)
	if err != nil {
		return err
	}

	// run in new thread so we can track progress asynchronously
	go func() {
		beganAt := time.Now()
		// make sure task status always correct even if it panicked
		defer func() {
			_, _ = runningTasks.Remove(task.ID)
			close(progress)
			if r := recover(); r != nil {
				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("run task failed: %v", r)
				}
			}
			finishedAt := time.Now()
			spentSeconds := finishedAt.Unix() - beganAt.Unix()
			if err != nil {
				err = models.Db.Model(task).Updates(map[string]interface{}{
					"status":        models.TASK_FAILED,
					"message":       err.Error(),
					"finished_at":   finishedAt,
					"spent_seconds": spentSeconds,
				}).Error
			} else {
				err = models.Db.Model(task).Updates(map[string]interface{}{
					"status":        models.TASK_COMPLETED,
					"message":       "",
					"finished_at":   finishedAt,
					"spent_seconds": spentSeconds,
				}).Error
			}
		}()
		// start execution
		logger.Info("start executing task ", task.ID)
		err = models.Db.Model(task).Updates(map[string]interface{}{
			"status":   models.TASK_RUNNING,
			"message":  "",
			"began_at": beganAt,
		}).Error
		if err != nil {
			logger.Error("update task state failed", err)
			return
		}
		err = plugins.RunPlugin(task.Plugin, options, progress, ctx)
	}()

	// read progress from working thread and save into database
	for p := range progress {
		logger.Info("running plugin progress", fmt.Sprintf(" %d %s %f%%", task.ID, task.Plugin, p*100))
		err = models.Db.Model(task).Updates(map[string]interface{}{
			"progress": p,
		}).Error
		if err != nil {
			logger.Error("save task progress failed", err)
			return err
		}
	}

	return err
}

func CancelTask(taskId uint64) error {
	cancel, err := runningTasks.Remove(taskId)
	if err != nil {
		return err
	}
	cancel()
	return nil
}
