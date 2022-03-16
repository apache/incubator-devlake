package core

import (
	"context"

	"gorm.io/gorm"
)

// prepare for temporal

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

type SubTaskContext interface {
	ExecContext
	TaskContext() TaskContext
}

// Transient resources needed for task execution
type TaskContext interface {
	ExecContext
	SetData(data interface{})
	SubTaskContext(subtask string) (SubTaskContext, error)
}

type SubTaskEntryPoint func(c SubTaskContext) error

type SubTaskMeta struct {
	Name             string
	EntryPoint       SubTaskEntryPoint
	EnabledByDefault bool
	Description      string
}

// implement this interface to let framework run tasks for you
type ManagedSubTasks interface {
	// return all available subtasks, framework will run them for you in order
	SubTaskMetas() []SubTaskMeta
	// based on task context and user input options, return data that shared among all subtasks
	PrepareTaskData(taskCtx TaskContext, options map[string]interface{}) (interface{}, error)
}
