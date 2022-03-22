package core

import (
	"context"

	"gorm.io/gorm"
)

type ProgressType int

const (
	TaskSetProgress ProgressType = iota
	TaskIncProgress
	SubTaskSetProgress
	SubTaskIncProgress
	SetCurrentSubTask
)

type RunningProgress struct {
	Type          ProgressType
	Current       int
	Total         int
	SubTaskName   string
	SubTaskNumber int
}

// This interface define all resources that needed for task/subtask execution
type ExecContext interface {
	GetName() string
	GetConfig(name string) string
	GetDb() *gorm.DB
	GetContext() context.Context
	GetLogger() Logger
	GetData() interface{}
	SetProgress(current int, total int)
	IncProgress(quantity int)
}

// This interface define all resources that needed for subtask execution
type SubTaskContext interface {
	ExecContext
	TaskContext() TaskContext
}

// This interface define all resources that needed for task execution
type TaskContext interface {
	ExecContext
	SetData(data interface{})
	SubTaskContext(subtask string) (SubTaskContext, error)
}

type SubTask interface {
	Execute() error
}

// All subtasks from plugins should comply to this prototype, so they could be orchestrated by framework
type SubTaskEntryPoint func(c SubTaskContext) error

// Meta data of a subtask
type SubTaskMeta struct {
	Name       string
	EntryPoint SubTaskEntryPoint
	// Required SubTask will be executed no matter what
	Required         bool
	EnabledByDefault bool
	Description      string
}

// Implement this interface to let framework run tasks for you
type PluginTask interface {
	// return all available subtasks, framework will run them for you in order
	SubTaskMetas() []SubTaskMeta
	// based on task context and user input options, return data that shared among all subtasks
	PrepareTaskData(taskCtx TaskContext, options map[string]interface{}) (interface{}, error)
}
