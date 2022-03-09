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
	SubTaskContext(subtask string) (SubTaskContext, error)
}

type SubTaskEntryPoint func(c SubTaskContext) error
